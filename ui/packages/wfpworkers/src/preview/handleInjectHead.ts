class ElementHandler {
  constructor() {}

  element(element) {
    // An incoming element, such as `div`
    element.append(
      `<meta
    property="og:title"
    content="Walnuts"
/>`,
      {html: true},
    );
    element.append(
      `<meta
      property="og:site_name"
      content="walnuts"
  />`,
      {html: true},
    );
    element.append(
      `<meta
        property="og:description"
        content=""
    />`,
      {html: true},
    );
  }

  comments(comment) {
    // An incoming comment
  }

  text(text) {
    // An incoming piece of text
  }
}

export async function handleInjectHead(request): Promise<Response> {
  const res = await fetch(request);
  const url = new URL(request.url);
  const segments = url.pathname.split('/');
  const key = segments.pop() || segments.pop();
  console.log('received share request: ' + key);

  return new HTMLRewriter().on('head', new ElementHandler()).transform(res);
}

export async function handleInjectHeadDB(request): Promise<Response> {
  const res = await fetch(request);
  const url = new URL(request.url);
  const segments = url.pathname.split('/');
  const key = segments.pop() || segments.pop();
  console.log('received share request: ' + key);

  return new HTMLRewriter().on('head', new ElementHandler()).transform(res);
}
