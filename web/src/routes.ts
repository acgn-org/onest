import { createBrowserRouter, type RouteObject } from "react-router";

import App from "./App";

import Downloads from "@page/Downloads";
import LogStream from "@page/LogStream";
import Items from "@page/Items";

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
        Component: Items,
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
