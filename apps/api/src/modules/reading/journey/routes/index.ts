import { JourneyRoutesV1 } from "./journey.v1.route.ts";

export const createJourneyRoutesV1 = () => {
  return new JourneyRoutesV1().router;
};
