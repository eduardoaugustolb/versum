import { DiscoveryRoutesV1 } from "./discovery.v1.route.ts";

export const createDiscoveryRoutesV1 = () => {
  return new DiscoveryRoutesV1().router;
};
