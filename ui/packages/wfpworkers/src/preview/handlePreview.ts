import {IRequest} from 'itty-router';

export async function handlePreview(
  request: IRequest,
  event: FetchEvent,
): Promise<Response> {
  let {params} = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }
  const key = params.key.replace('.png', '');

  if (key === '') {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request',
    });
  }

  console.log(key);

  const cacheUrl = new URL(request.url);
  const cacheKey = new Request(cacheUrl.toString(), request);
  console.log(`checking for cache key: ${cacheUrl}`);
  const cache = caches.default;
  let response = await cache.match(cacheKey);

  if (!response) {
    console.log(
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`,
    );

    const resp = await fetch(
      new Request(PREVIEW_ENDPOINT + '/generate/sh/' + key),
      {
        headers: {
          'X-CUSTOM-AUTH-KEY': AUTH_KEY,
        },
        cf: {
          cacheTtl: 60 * 24 * 60 * 60,
          cacheEverything: true,
        },
      },
    );

    response = new Response(resp.body, resp);

    //don't cache errors
    if (resp.ok) {
      response.headers.set('Cache-Control', 'max-age=5184000');
      event.waitUntil(cache.put(cacheKey, response.clone()));
    }
  }

  return response;
}
