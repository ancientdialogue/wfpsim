import {
  Button,
  ButtonGroup,
  Intent,
  Position,
  Toaster,
} from '@blueprintjs/core';
import {Executor, ExecutorSupplier} from '@gcsim/executors';
import {SimResults} from '@gcsim/types';
import {ResultSource} from '@gcsim/ui/src/Pages';
import LoadingToast from '@gcsim/ui/src/Pages/Viewer/Components/LoadingToast';
import Warnings from '@gcsim/ui/src/Pages/Viewer/Components/Warnings';
import Results from '@gcsim/ui/src/Pages/Viewer/Tabs/Results';
import classNames from 'classnames';
import React, {useCallback, useMemo} from 'react';
import {useRunningState} from '../../UI';
import {Share} from './Share';

type ViewerProps = {
  data: SimResults;
  hash: string | null;
  exec: ExecutorSupplier<Executor>;
};

export function WebViewer({data}: {data: SimResults}) {
  const names = useMemo(
    () => data?.character_details?.map((c) => c.name),
    [data?.character_details],
  );
  const copyToast = React.useRef<Toaster>(null);
  const copy = () => {
    navigator.clipboard.writeText(data.config_file ?? '').then(() => {
      copyToast.current?.show({
        message: 'Link copied to clipboard!',
        intent: Intent.SUCCESS,
        timeout: 2000,
      });
    });
  };
  const home = () => {
    window.location.replace('/');
  };
  return (
    <div className="flex flex-col flex-grow w-full pb-6">
      <div className="ml-auto">
        <ButtonGroup>
          <Button onClick={copy} intent="primary">
            Copy Config
          </Button>
          <Button onClick={home} intent="success">
            Simulator
          </Button>
        </ButtonGroup>
      </div>
      <Warnings data={data} />
      <div className="basis-full pt-0 mt-1">
        <Results data={data} running={false} names={names} />,
      </div>
      <Toaster ref={copyToast} position={Position.TOP} />
    </div>
  );
}

export function Viewer({data, hash = '', exec}: ViewerProps) {
  const running = useRunningState(exec);
  const names = useMemo(
    () => data?.character_details?.map((c) => c.name),
    [data?.character_details],
  );
  const cancel = useCallback(() => exec().cancel(), [exec]);
  return (
    <div className="flex flex-col flex-grow w-full pb-6">
      <Warnings data={data} />
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <ViewerNav data={data} running={running} />
      </div>
      <div className="basis-full pt-0 mt-1">
        <Results data={data} running={false} names={names} />,
      </div>
      <LoadingToast
        cancel={cancel}
        running={running}
        src={ResultSource.Generated}
        error={null}
        current={data?.statistics?.iterations}
        total={data?.simulator_settings?.iterations}
      />
    </div>
  );
}

type NavProps = {
  data: SimResults | null;
  running: boolean;
};

const btnClass = classNames('hidden ml-[7px] sm:flex');

function ViewerNav({data, running}: NavProps) {
  const copyToast = React.useRef<Toaster>(null);
  const shareState = React.useState<string | null>(null);

  return (
    <div>
      <Share
        copyToast={copyToast}
        shareState={shareState}
        data={data}
        running={running}
        className={btnClass}
      />
      <Toaster ref={copyToast} position={Position.TOP_RIGHT} />
    </div>
  );
}
