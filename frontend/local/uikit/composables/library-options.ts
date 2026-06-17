import { inject } from "vue";

import type { ITNLibraryOptions } from "../interfaces";

export const useLibraryOptions = () => {
  return inject<Readonly<ITNLibraryOptions>>("libraryOptions", {});
};
