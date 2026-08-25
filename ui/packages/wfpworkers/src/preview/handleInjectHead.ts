const escapeMetaValue = (value = '') =>
  value
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');

async function getShareDescription(host, key, prefix) {
  const sharePath = prefix === 'db' ? `db/${key}` : key;
  const shareUrl = new URL('/api/share/' + sharePath, host + '/');

  try {
    const res = await fetch(shareUrl);
    if (!res.ok) {
      return '';
    }

    const data = await res.json();
    const names = (data.character_details ?? [])
      .map((character) => character?.name)
      .filter((name) => typeof name === 'string' && name.trim().length > 0)
      .slice(0, 4);

    return names.join(', ');
  } catch {
    return '';
  }
}

class ElementHandler {
  private key;
  private host;
  private prefix;
  private description;

  constructor(host, key, prefix, description) {
    this.key = key;
    this.host = host;
    this.prefix = prefix;
    this.description = description;
  }

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
        content="${escapeMetaValue(this.description)}"
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
  const host = url.protocol + '//' + url.host;
  const description = await getShareDescription(host, key, '');
  console.log('received share request: ' + key);

  return new HTMLRewriter()
    .on('head', new ElementHandler(host, key, '', description))
    .transform(res);
}

export async function handleInjectHeadDB(request): Promise<Response> {
  const res = await fetch(request);
  const url = new URL(request.url);
  const segments = url.pathname.split('/');
  const key = segments.pop() || segments.pop();
  const host = url.protocol + '//' + url.host;
  const description = await getShareDescription(host, key, 'db');
  console.log('received share request: ' + key);

  return new HTMLRewriter()
    .on('head', new ElementHandler(host, key, 'db', description))
    .transform(res);
}
