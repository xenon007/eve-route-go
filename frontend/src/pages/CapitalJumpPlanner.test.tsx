import React from "react";
import { render } from "@testing-library/react";
import { ThemeProvider } from "@material-ui/core/styles";
import theme from "../theme";
import CapitalJumpPlanner from "./CapitalJumpPlanner";

test("CapitalJumpPlanner snapshot", () => {
  const { container } = render(
    <ThemeProvider theme={theme}>
      <CapitalJumpPlanner />
    </ThemeProvider>,
  );
  expect(container).toMatchSnapshot();
});
