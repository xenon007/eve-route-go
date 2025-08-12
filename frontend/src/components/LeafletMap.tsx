import React, { useEffect, useRef } from "react";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import { CapitalSystem } from "../response";

// Props задаёт список систем для отрисовки маршрута.
type Props = {
  systems: CapitalSystem[];
};

/**
 * LeafletMap отображает маршрут капитальных прыжков на карте Leaflet.
 */
export default function LeafletMap({ systems }: Props) {
  const mapRef = useRef<HTMLDivElement>(null);
  const mapInst = useRef<L.Map>();

  // Инициализация карты один раз после монтирования.
  useEffect(() => {
    if (!mapRef.current || mapInst.current) {
      return;
    }
    mapInst.current = L.map(mapRef.current).setView([0, 0], 1);
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "&copy; OpenStreetMap contributors",
    }).addTo(mapInst.current);
  }, []);

  // Отрисовка линии маршрута при обновлении систем.
  useEffect(() => {
    const map = mapInst.current;
    if (!map) {
      return;
    }
    // Удаляем предыдущие линии.
    map.eachLayer((layer) => {
      if (layer instanceof L.Polyline && !(layer instanceof L.TileLayer)) {
        map.removeLayer(layer);
      }
    });
    if (systems.length === 0) {
      return;
    }
    console.info("Drawing route with", systems.length, "systems");
    const scale = 1e16; // перевод метров в условные градусы
    const coords = systems.map(
      (s) => [s.y / scale, s.x / scale] as [number, number],
    );
    const polyline = L.polyline(coords, { color: "red" }).addTo(map);
    map.fitBounds(polyline.getBounds());
  }, [systems]);

  return <div ref={mapRef} style={{ height: "400px" }} />;
}
