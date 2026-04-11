import { DiscoveryRoutes } from "./discovery.route.ts";

export const createDiscoveryRoutes = () => {
  return new DiscoveryRoutes().router;
};
