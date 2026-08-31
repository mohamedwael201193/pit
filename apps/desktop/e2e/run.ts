import { assertOfficialLinks } from "./links.spec.ts";
import { assertHonesty } from "./honesty.spec.ts";
import { assertNamedResearchWhy } from "./explain.spec.ts";
import { assertDesktopCopy } from "./onboarding.spec.ts";
import { assertAuthorizeClosed } from "./authorize.spec.ts";
import { assertLedgerUnsignedCopy } from "./ledger.spec.ts";
import { assertPreviewMutationCopy } from "./preview.spec.ts";
import { assertSessionCopy } from "./session.spec.ts";
import { assertNetworkToggle } from "./network.spec.ts";
import { assertReducedMotionCopy } from "./motion.spec.ts";
import { assertTeeNeverGreenIdle, assertDirectNotGreenWithoutAuth } from "./readiness.spec.ts";
import { assertNamedErrors } from "./errors.spec.ts";
import { assertCapitalFloorIgnoresDust, assertDeskRanksExecutableFirst, assertPolicyClipTightCopy, assertAgenticIdPartial, assertDeskHeadlinePrefersExecutable } from "./capital.spec.ts";
import { assertShellFilters } from "./shell.spec.ts";
import { assertChatAgentCopy } from "./chat-agent.spec.ts";
import { assertLiveAgentPipeline } from "./pipeline.spec.ts";
import {
  assertOnboardPairFirst,
  assertOnboardDoesNotInventReady,
  assertOnboardReadyRequiresAllGates,
  assertNextFixPairsFirst,
  assertNextFixDoesNotSkipPairWhenSessionExists,
  assertOnboardInputFromDoctor,
} from "./onboard.spec.ts";

assertOfficialLinks();
assertHonesty();
assertNamedResearchWhy();
assertDesktopCopy();
assertAuthorizeClosed("AUTHORIZE", true);
assertAuthorizeClosed("yes", true);
assertAuthorizeClosed("AUTHORIZE", false);
assertLedgerUnsignedCopy("unsigned venue payload");
assertPreviewMutationCopy("mutation invalidates the preview");
assertSessionCopy("PIT never prints the session key");
assertNetworkToggle("mainnet");
assertNetworkToggle("testnet");
assertReducedMotionCopy();
assertTeeNeverGreenIdle();
assertDirectNotGreenWithoutAuth();
assertNamedErrors();
assertCapitalFloorIgnoresDust();
assertDeskRanksExecutableFirst();
assertPolicyClipTightCopy();
assertAgenticIdPartial();
assertDeskHeadlinePrefersExecutable();
assertShellFilters();
assertChatAgentCopy();
assertLiveAgentPipeline();
assertOnboardPairFirst();
assertOnboardDoesNotInventReady();
assertOnboardReadyRequiresAllGates();
assertNextFixPairsFirst();
assertNextFixDoesNotSkipPairWhenSessionExists();
assertOnboardInputFromDoctor();
console.log("desktop e2e copy harness ok");
