import {Button, ButtonGroup, Callout, Intent} from '@blueprintjs/core';
import {Executor, ExecutorSupplier} from '@gcsim/executors';
import {ActionList} from '@gcsim/ui/src/Pages/Simulator/Components';
import {debounce} from 'lodash-es';
import React from 'react';

export function Simulator({
  exec,
  run,
  openSettings,
}: {
  exec: ExecutorSupplier<Executor>;
  run: (string) => void;
  openSettings?: () => void;
}) {
  const [cfg, setCfg] = React.useState<string>(() => {
    const localData = localStorage.getItem('cfg');
    return localData ? localData : '';
  });
  React.useEffect(() => {
    localStorage.setItem('cfg', cfg);
  }, [cfg]);

  // check worker ready state every 250ms so run button becomes available when workers do
  const [isReady, setReady] = React.useState<boolean | null>(null);
  React.useEffect(() => {
    const interval = setInterval(() => {
      exec()
        .ready()
        .then((res) => setReady(res));
    }, 250);
    return () => clearInterval(interval);
  }, [exec]);
  const [err, setErr] = React.useState('');

  const validated = useConfigValidateListener(exec, cfg, isReady, setErr);
  const onChange = (next: string) => {
    setCfg(next);
  };
  return (
    <div className="px-2 py-4 w-full 2xl:mx-auto 2xl:container flex flex-col gap-2">
      <div className="mt-2 mb-2">
        gcsim devs refuse to implement unreleased characters, so we did it for
        them. The following are characters/weapons/artifacts currently not
        implemented in gcsim that we have added here:
        <ul className="list-disc pl-4">
          <li>Xilonen</li>
          <li>Mualani</li>
        </ul>
        <p className=" font-bold">
          See{' '}
          <a
            href="https://gist.github.com/ac1dgit/89c92d666fb5805266797af635d56464"
            target="_blank"
            rel="noopener noreferrer">
            here
          </a>{' '}
          for all implementation assumptions for pre-release characters
        </p>
      </div>
      <ActionList cfg={cfg} onChange={onChange} />
      <div className="sticky bottom-0 flex flex-col gap-y-1 bg-[#450a0a]">
        {err !== '' && cfg !== '' ? (
          <div className="pl-2 pr-2 pt-2 mt-1">
            <Callout intent={Intent.DANGER} title="Error: Config Invalid">
              <pre className="whitespace-pre-wrap pl-5">{err}</pre>
            </Callout>
          </div>
        ) : null}
        <div className="p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap place-items-center gap-x-1 gap-y-1">
          <ButtonGroup className="basis-full wide:basis-2/3 p-1 flex flex-row flex-wrap">
            {openSettings === undefined ? null : (
              <Button
                className="!basis-full md:!basis-1/2"
                icon="settings"
                onClick={openSettings}
                text={'Settings'}
              />
            )}
            <Button
              icon="play"
              intent="primary"
              className="!basis-full md:!basis-1/2"
              onClick={() => run(cfg)}
              loading={!isReady}
              disabled={err !== '' && !validated}
              text={'Run'}
            />
          </ButtonGroup>
        </div>
      </div>
    </div>
  );
}

export function useConfigValidateListener(
  exec: ExecutorSupplier<Executor>,
  cfg: string,
  isReady: boolean | null,
  setErr: (str: string) => void,
): boolean {
  const [validated, setValidated] = React.useState(false);
  const debounced = React.useRef(debounce((x: () => void) => x(), 200));

  React.useEffect(() => {
    if (!isReady) {
      return;
    }

    if (cfg === '') {
      return;
    }

    setValidated(false);
    debounced.current(() => {
      exec()
        .validate(cfg)
        .then(
          (res) => {
            console.log('all is good');
            setErr('');
            //check if there are any warning msgs
            if (res.errors) {
              let msg = '';
              res.errors.forEach((err) => {
                msg += err + '\n';
              });
              setErr(msg);
            }
            setValidated(true);
          },
          (err) => {
            //set error state
            setErr(err);
            setValidated(false);
          },
        );
    });
  }, [exec, cfg, setErr, isReady]);

  return validated;
}
