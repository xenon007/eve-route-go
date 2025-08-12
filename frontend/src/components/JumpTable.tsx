import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from "@material-ui/core";
import { useTranslation } from "react-i18next";

const FUEL_PER_LY = 1000; // количество изотопов на один световой год
const LY_IN_METERS = 9.4607e15;

interface System {
  id: number;
  name: string;
  x: number;
  y: number;
  z: number;
}

/**
 * JumpTable отображает список прыжков с расстоянием и расходом топлива.
 */
export default function JumpTable({ systems }: { systems: System[] }) {
  const { t } = useTranslation();
  if (systems.length < 2) {
    return null;
  }

  const rows = systems.slice(1).map((s, i) => {
    const prev = systems[i];
    const dx = prev.x - s.x;
    const dy = prev.y - s.y;
    const dz = prev.z - s.z;
    const distance = Math.sqrt(dx * dx + dy * dy + dz * dz) / LY_IN_METERS;
    return {
      system: s.name,
      distance,
      fuel: distance * FUEL_PER_LY,
    };
  });

  const totalFuel = rows.reduce((sum, r) => sum + r.fuel, 0);

  return (
    <>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>{t("capital.system")}</TableCell>
            <TableCell>{t("capital.distance")}</TableCell>
            <TableCell>{t("capital.fuel")}</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((r, i) => (
            <TableRow key={i}>
              <TableCell>{r.system}</TableCell>
              <TableCell>{r.distance.toFixed(2)}</TableCell>
              <TableCell>{Math.round(r.fuel)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Typography variant="caption">
        {t("capital.total-fuel")}: {Math.round(totalFuel)}
      </Typography>
    </>
  );
}
