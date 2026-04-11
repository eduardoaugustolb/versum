import prompts from "prompts";

export const createConfirmActionPrompt = async (action: string) => {
  const key = `confirm_${action.replaceAll(" ", "_").toLowerCase()}`;
  const reply = await prompts({
    type: "confirm",
    name: key,
    message: action,
  });
  return reply[key];
};
