import { JourneyRoutes } from "./journey.route.ts";

export const createJourneyRoutes = () => {
  return new JourneyRoutes().router;
};
