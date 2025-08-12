import React from "react";
import { render, fireEvent, screen, waitFor } from "@testing-library/react";
import { ThemeProvider } from "@material-ui/core/styles";
import theme from "../theme";
import Capital from "./Capital";
import axios from "axios";

jest.mock("axios");

// Заменяем сложные компоненты простыми моками для изоляции теста.
jest.mock(
  "../components/SystemInput",
  () =>
    ({ fieldId, onChange }: { fieldId: string; onChange: Function }) => (
      <input
        data-testid={fieldId}
        onChange={(e) => onChange((e.target as HTMLInputElement).value)}
      />
    ),
);

jest.mock(
  "../components/LeafletMap",
  () =>
    ({ systems }: { systems: any[] }) => (
      <div data-testid="map">{systems.length}</div>
    ),
);

test("Capital snapshot", () => {
  const { container } = render(
    <ThemeProvider theme={theme}>
      <Capital />
    </ThemeProvider>,
  );
  expect(container).toMatchSnapshot();
});

test("renders route with region and fuel", async () => {
  (axios.get as jest.Mock).mockResolvedValueOnce({
    data: {
      route: [
        { id: 1, name: "Start", regionName: "Alpha", x: 0, y: 0, z: 0 },
        {
          id: 2,
          name: "End",
          regionName: "Beta",
          x: 9.4607e15,
          y: 0,
          z: 0,
        },
      ],
    },
  });

  render(
    <ThemeProvider theme={theme}>
      <Capital />
    </ThemeProvider>,
  );

  fireEvent.change(screen.getByTestId("start-system"), {
    target: { value: "Start" },
  });
  fireEvent.change(screen.getByTestId("end-system"), {
    target: { value: "End" },
  });
  const btn = document.getElementById("find-route") as HTMLButtonElement;
  fireEvent.click(btn);

  await waitFor(() => screen.getByText("End"));
  expect(screen.getByText("Beta")).toBeInTheDocument();
  expect(screen.getByText("1000")).toBeInTheDocument();
  expect(screen.getByTestId("map")).toHaveTextContent("2");
});
