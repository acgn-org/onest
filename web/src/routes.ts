import { createBrowserRouter, type RouteObject } from "react-router";

import App from "./App";

import Downloads from "@page/Downloads";
import LogStream from "@page/LogStream";

const routes: RouteObject[] = [
  {
    path: "/",
    Component: App,
    children: [
      {
        index: true,
        Component: Downloads,
      },
      {
        path: "/items",
        element: "Items",
      },
      {
        path: "/time-machine",
        element: "Time Machine",
      },
      {
        path: "/log-stream",
        Component: LogStream,
      },
      {
        path: "*",
        element: "NotFound",
      },
    ],
  },
];

export const router = createBrowserRouter(routes);
export default router;
