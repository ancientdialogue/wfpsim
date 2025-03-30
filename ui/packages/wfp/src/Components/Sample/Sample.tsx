import {Executor, ExecutorSupplier} from '@gcsim/executors';
import {SimResults} from '@gcsim/types';
import SampleUI, {useSample} from '@gcsim/ui/src/Pages/Viewer/Tabs/Sample';
import {useCallback} from 'react';
import {useRunningState} from '../../UI';

type SampleProps = {
  data: SimResults;
  exec: ExecutorSupplier<Executor>;
};

export function Sample({data, exec}: SampleProps) {
  const sampler = useCallback(
    (cfg: string, seed: string) => exec().sample(cfg, seed),
    [exec],
  );
  const running = useRunningState(exec);
  const sample = useSample(running, data, true, sampler);
  return (
    <div className="flex flex-col flex-grow w-full pb-6">
      <SampleUI
        sampler={sampler}
        data={data}
        sample={sample}
        running={running}
      />
    </div>
  );
}
