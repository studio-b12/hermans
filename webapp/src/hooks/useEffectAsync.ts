import { DependencyList, useEffect } from "react";

const useEffectAsync = (
  effect: () => Promise<void>,
  deps: DependencyList,
  destructor?: () => void,
) => {
  useEffect(() => {
    effect();
    return destructor;
  }, deps);
};

export default useEffectAsync;
