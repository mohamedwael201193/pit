export function refuseArm(): never {
  throw new Error("mission_arm_denied");
}

export const canArm = false;
