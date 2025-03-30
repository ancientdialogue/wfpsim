import {Button, Classes, Dialog, Tab, Tabs} from '@blueprintjs/core';
import {Executor, ExecutorSupplier} from '@gcsim/executors';
import {SimResults} from '@gcsim/types';
import axios from 'axios';
import {throttle} from 'lodash-es';
import {ReactNode, useEffect, useState} from 'react';
import {Sample} from './Components/Sample/Sample';
import {Simulator} from './Components/Simulator/Simulator';
import {Viewer, WebViewer} from './Components/Viewer/Viewer';

type UIProps = {
  exec: ExecutorSupplier<Executor>;
  children: ReactNode;
};

type simRunResult = {
  data: SimResults;
  hash: string;
} | null;

export const UI = ({exec, children}: UIProps) => {
  const [tabId, setTabId] = useState('simulator');
  const [results, setResults] = useState<simRunResult>(null);
  const [webResult, setWebResult] = useState<SimResults | null>(null);
  const [err, setError] = useState<string>('');
  const [settingsOpen, setSettingsOpen] = useState<boolean>(false);

  let key = '';
  const res =
    /sh\/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$/.exec(
      window.location.pathname,
    );
  if (res && res.length >= 1) {
    key = res[1];
  }

  useEffect(() => {
    if (key !== '') {
      axios
        .get('/api/share/' + key, {timeout: 30000})
        .then((resp) => {
          setWebResult(resp.data);
          console.log(resp.data);
        })
        .catch((e) => {
          setError(e.message);
        });
    } else {
      setWebResult(null);
    }
  }, [key, setResults]);

  if (key != '') {
    if (webResult === null) {
      return <div>Loading share please wait</div>;
    }

    return (
      <div className="flex flex-col flex-grow w-full pb-6">
        <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
          <WebViewer data={webResult} />
        </div>
      </div>
    );
  }

  const runSim = (cfg: string) => {
    console.log('starting run');
    setResults(null);
    setError('');

    const updateResult = throttle(
      (res: SimResults, hash: string) => {
        if (tabId !== 'results') {
          setTabId('results');
        }
        setResults({data: res, hash: hash});
      },
      100,
      {leading: true, trailing: true},
    );

    exec()
      .run(cfg, updateResult)
      .catch((err) => {
        console.log('problems :(', err);
        setError(err);
      });
  };

  const tabs: {[k: string]: React.ReactNode} = {
    simulator: (
      <Simulator
        exec={exec}
        run={runSim}
        openSettings={() => setSettingsOpen(true)}
      />
    ),
    results:
      results === null ? (
        <></>
      ) : (
        <Viewer data={results.data} hash={results.hash} exec={exec} />
      ),
    sample:
      results === null ? <></> : <Sample data={results.data} exec={exec} />,
  };

  if (err !== '') {
    return (
      <div>
        oops something went wrong:
        <br />
        {JSON.stringify(err)}
        <br />
        <Button
          icon="refresh"
          onClick={() => {
            exec().cancel();
            setError('');
            setResults(null);
          }}
          intent="primary">
          Reload
        </Button>
      </div>
    );
  }

  return (
    <div className="flex flex-col flex-grow w-full pb-6">
      <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container">
        <Tabs selectedTabId={tabId} onChange={(s) => setTabId(s as string)}>
          <Tab
            id="simulator"
            className="focus:outline-none"
            title="Simulator"></Tab>
          <Tab
            id="results"
            className="focus:outline-none"
            title="Results"
            disabled={results === null}></Tab>
          <Tab
            id="sample"
            className="focus:outline-none"
            title="Sample"
            disabled={results === null}></Tab>
          <Tabs.Expander />
        </Tabs>
      </div>
      {tabs[tabId]}
      <ExecutorSettings
        isOpen={settingsOpen}
        onClose={() => setSettingsOpen(false)}>
        {children}
      </ExecutorSettings>
    </div>
  );
};

export function useRunningState(exec: ExecutorSupplier<Executor>): boolean {
  const [isRunning, setRunning] = useState(true);

  useEffect(() => {
    const check = setInterval(() => {
      setRunning(exec().running());
    }, 100 - 50);
    return () => clearInterval(check);
  }, [exec]);

  return isRunning;
}

type ExecutorSettingsProps = {
  children: ReactNode;
  isOpen: boolean;
  onClose: () => void;
};

const ExecutorSettings = ({
  children,
  isOpen,
  onClose,
}: ExecutorSettingsProps) => {
  return (
    <Dialog
      isOpen={isOpen}
      //   onClose={onClose}
      title="Settings"
      icon="settings"
      className="!pb-0">
      <div className={Classes.DIALOG_BODY}>{children}</div>
      <div className={Classes.DIALOG_FOOTER}>
        <Button onClick={onClose}>Close</Button>
      </div>
    </Dialog>
  );
};
