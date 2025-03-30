import { IRequest } from "itty-router";
import pako from "pako";

export async function handleView(
  request: IRequest,
  event: FetchEvent
): Promise<Response> {
  let { params } = request;
  if (!params || !params.key) {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
    });
  }

  const key = params.key;

  if (key === "") {
    return new Response(null, {
      status: 400,
      statusText: "Bad Request",
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
      `Response for request url: ${request.url} not present in cache. Fetching and caching request.`
    );

    //try grabbing from kv
    try {
      const compressed: unknown = await WFPSIM_KV.get(key, {
        type: "arrayBuffer",
      });
      if (compressed === null) {
        return new Response("share not found", { status: 404 });
      }
      const data = pako.inflate(compressed as ArrayBuffer);
      response = new Response(data, response);
      response.headers.append("Cache-Control", "s-maxage=1800");
      response.headers.append("Content-Encoding", "gzip");

      event.waitUntil(cache.put(cacheKey, response.clone()));
    } catch {
      return new Response("getting from kv failed", { status: 500 });
    }
  } else {
    console.log(`cache hit for: ${request.url}`);
  }

  return response;
}
