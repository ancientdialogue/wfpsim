import {
  Button,
  Classes,
  Dialog,
  Icon,
  InputGroup,
  Intent,
  Label,
  NonIdealState,
  Spinner,
  SpinnerSize,
  Toaster,
} from "@blueprintjs/core";
import { SimResults } from "@gcsim/types";
import axios from "axios";
import classNames from "classnames";
import React, { RefObject } from "react";

export const Share = ({
  running,
  copyToast,
  data,
  className,
  shareState,
}: ShareProps) => {
  const [isOpen, setOpen] = React.useState(false);
  const [shareLink, setShareLink] = shareState;

  // change the set link if url changes or rerun
  React.useEffect(() => {
    setShareLink(extractFromLocation(window.location.pathname));
  }, [window.location.pathname, setShareLink, data?.config_file]);

  const handleShare = () => {
    if (data === null || shareLink != null) {
      return;
    }

    axios
      .post("/api/share", data)
      .then((resp) => {
        setShareLink(link("sh", resp.data));
      })
      .catch((err) => {
        console.log(err);
      });
  };

  const copy = () => {
    navigator.clipboard.writeText(shareLink ?? "").then(() => {
      copyToast.current?.show({
        message: "Link copied to clipboard!",
        intent: Intent.SUCCESS,
        timeout: 2000,
      });
    });
  };

  return (
    <>
      <Button
        icon={<Icon icon="link" className="!mr-0" />}
        intent={Intent.PRIMARY}
        disabled={running || data == null}
        onClick={() => {
          handleShare();
          setOpen(true);
        }}
      >
        <div className={className}>Share</div>
      </Button>
      <Dialog
        isOpen={isOpen}
        onClose={() => setOpen(false)}
        title={"Create share link"}
        icon="link"
        className="!pb-0"
      >
        <div
          className={classNames(
            Classes.DIALOG_BODY,
            "flex flex-col justify-center gap-2"
          )}
        >
          <DialogBody shareLink={shareLink} copy={copy} />
        </div>
      </Dialog>
    </>
  );
};

type DialogProps = {
  shareLink: string | null;
  copy: () => void;
};

const DialogBody = ({ shareLink, copy }: DialogProps) => {
  if (shareLink == null) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  return (
    <Label>
      Share Link
      <InputGroup
        readOnly={true}
        fill={true}
        onFocus={(e) => {
          e.target.select();
          copy();
        }}
        value={shareLink ?? ""}
        className={classNames({ "bp4-skeleton": shareLink == null })}
        large={true}
        rightElement={<Button icon="duplicate" onClick={() => copy()} />}
      />
    </Label>
  );
};

export type ShareProps = {
  running: boolean;
  copyToast: RefObject<Toaster>;
  data: SimResults | null;
  shareState: [string | null, (link: string | null) => void];
  className?: string;
};

export function link(route: string, id: string): string {
  return `${window.location.protocol}//${window.location.host}/${route}/${id}`;
}

export function extractFromLocation(location: string) {
  if (location.startsWith("/sh/")) {
    return link("sh", location.substring(location.lastIndexOf("/") + 1));
  } else if (location.startsWith("/db/")) {
    return link("db", location.substring(location.lastIndexOf("/") + 1));
  }
  return null;
}
