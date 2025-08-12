import React, { useState } from "react";
import { Box, Button, Grid, Typography } from "@material-ui/core";
import axios from "axios";
import SystemInput from "../components/SystemInput";

interface ApiResponse {
  route: Array<{ id: number; name: string }>;
}

/**
 * Страница планировщика прыжков капитальных кораблей.
 * Позволяет выбрать начальную и конечную системы
 * и рассчитать маршрут через API.
 */
export default function CapitalJumpPlanner() {
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");
  const [names, setNames] = useState<string[]>([]);
  const [message, setMessage] = useState("");

  const findRoute = () => {
    setMessage("");
    if (!start || !end) {
      return;
    }
    console.log("Requesting capital route", start, end);
    axios
      .get<ApiResponse>(`/api/capital?start=${start}&end=${end}`)
      .then((r) => {
        setNames(r.data.route.map((s) => s.name));
      })
      .catch((err) => {
        console.error("Failed to fetch capital route", err);
        setMessage("No route found");
      });
  };

  return (
    <Grid container spacing={2} className="card">
      <Grid item xs={12}>
        <Typography variant="h6" align="center">
          Capital Jump Planner
        </Typography>
      </Grid>

      <Grid item sm={5} xs={12}>
        <Box display="flex" justifyContent="center">
          <SystemInput
            fieldId="start-system"
            fieldName="Start"
            onChange={setStart}
            fieldValue={start}
            findRoute={findRoute}
          />
        </Box>
      </Grid>
      <Grid item sm={2} xs={12}></Grid>
      <Grid item sm={5} xs={12}>
        <Box display="flex" justifyContent="center">
          <SystemInput
            fieldId="end-system"
            fieldName="End"
            onChange={setEnd}
            fieldValue={end}
            findRoute={findRoute}
          />
        </Box>
      </Grid>

      <Grid item xs={12}>
        <Box display="flex" justifyContent="center">
          <Button
            variant="contained"
            color="primary"
            onClick={findRoute}
            disabled={!(start && end)}
          >
            Find Route
          </Button>
        </Box>
      </Grid>

      <Grid item xs={12}>
        {message && <Typography>{message}</Typography>}
        {!message && names.length > 0 && (
          <ol>
            {names.map((n, i) => (
              <li key={i}>{n}</li>
            ))}
          </ol>
        )}
      </Grid>
    </Grid>
  );
}
