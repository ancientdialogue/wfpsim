import { Validator } from "@cfworker/json-schema";
import SimulationResult from "./schema/SimulationResult.json";

export const validator = new Validator(SimulationResult);
