import { createBrowserRouter, type RouteObject } from "react-router";

import App from "./App";

const routes: RouteObject[] = [
  {
    path: "/",
    Component: App,
    children: [
      {
        index: true,
        element: "Downloads",
      },
      {
        path: "/time-machine",
        element: "Time Machine",
      },
      {
        path: "/log-stream",
        element: "Log Stream",
      },
    ],
  },
];

export const router = createBrowserRouter(routes);
export default router;
