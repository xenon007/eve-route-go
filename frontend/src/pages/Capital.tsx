import React, { useContext, useState } from "react";
import { Box, Button, Grid, Typography } from "@material-ui/core";
import { useTranslation } from "react-i18next";
import axios from "axios";
import SystemInput from "../components/SystemInput";
import LeafletMap from "../components/LeafletMap";
import { GlobalDataContext } from "../GlobalDataContext";
import { CapitalSystem, ResponseCapital } from "../response";

/**
 * Страница Capital Jump Planner.
 * Позволяет рассчитать маршрут прыжков капитального корабля
 * и отобразить его на карте.
 */
export default function Capital() {
  const { t } = useTranslation();
  const globalData = useContext(GlobalDataContext);
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");

  const [route, setRoute] = useState<CapitalSystem[]>([]);
  const [names, setNames] = useState<string[]>([]);
  const [message, setMessage] = useState("");

  const findRoute = () => {
    setMessage("");
    axios
      .get<ResponseCapital>(
        `${globalData.domain}/api/capital?start=${start}&end=${end}`,
      )
      .then((r) => {

        console.info("Received capital route", r.data.route);
        setRoute(r.data.route);
        setNames(r.data.route.map((s) => s.name));
      })
      .catch(() => {
        setMessage(t("capital.no-route"));
      });
  };

  return (
    <Grid container spacing={2} className="card">
      <Grid item xs={12}>
        <Typography variant="h6" align="center">
          {t("capital.title")}
        </Typography>
      </Grid>

      <Grid item sm={5} xs={12}>
        <Box display="flex" justifyContent="center">
          <SystemInput
            fieldId="start-system"
            fieldName={t("capital.start")}
            onChange={setStart}
            findRoute={findRoute}
            fieldValue={start}
          />
        </Box>
      </Grid>
      <Grid item sm={2} xs={12}></Grid>
      <Grid item sm={5} xs={12}>
        <Box display="flex" justifyContent="center">
          <SystemInput
            fieldId="end-system"
            fieldName={t("capital.end")}
            onChange={setEnd}
            findRoute={findRoute}
            fieldValue={end}
          />
        </Box>
      </Grid>

      <Grid item xs={12}>
        <Box display="flex" justifyContent="center">
          <Button
            id="find-route"
            variant="contained"
            color="primary"
            onClick={findRoute}
            disabled={!(start && end)}
          >
            {t("capital.find")}
          </Button>
        </Box>
      </Grid>

      <Grid item sm={4} xs={12}>
        {message && <Typography>{message}</Typography>}
        {!message && <JumpTable systems={systems} />}
      </Grid>
      <Grid item sm={8} xs={12}>
        <LeafletMap systems={route} />
      </Grid>
    </Grid>
  );
}
