import {FormGroup, InputGroup} from '@blueprintjs/core';
import {ExecutorSupplier, ServerExecutor} from '@gcsim/executors';
import React, {ReactNode, useRef} from 'react';
import {useTranslation} from 'react-i18next';
import {UI} from './UI';

let exec: ServerExecutor | undefined;
const urlKey = 'server-mode-url';
const defaultURL = 'http://127.0.0.1:54321';

const ServerMode = ({children}: {children: ReactNode}) => {
  const {t} = useTranslation();
  const [url, setURL] = React.useState<string>((): string => {
    const saved = localStorage.getItem(urlKey);
    if (saved === null) {
      localStorage.setItem(urlKey, defaultURL);
      return defaultURL;
    }
    return saved;
  });
  React.useEffect(() => {
    localStorage.setItem(urlKey, url);
  }, [url]);

  React.useEffect(() => {
    if (exec != null) {
      exec.set_url(url);
    }
  }, [url]);

  const supplier = useRef<ExecutorSupplier<ServerExecutor>>(() => {
    if (exec == null) {
      exec = new ServerExecutor(url);
    }
    return exec;
  });

  return (
    <UI exec={supplier.current}>
      <FormGroup className="!m-0" label={t('simple.workers')}>
        {children}
        <FormGroup
          helperText={
            t('simple.server_mode_default') + 'http://127.0.0.1:54321'
          }
          label={t('simple.server_mode_url')}
          labelFor="text-input"
          labelInfo={t('simple.server_mode_required')}>
          <InputGroup
            id="text-input"
            value={url}
            onChange={(e) => {
              setURL(e.target.value);
            }}
            fill
          />
        </FormGroup>
      </FormGroup>
    </UI>
  );
};

export default ServerMode;
