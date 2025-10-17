import {uuid} from '@cfworker/uuid';
import {IRequest} from 'itty-router';
import pako from 'pako';
import {validator} from './validation';

export async function handleShare(request: IRequest): Promise<Response> {
  let content: any;
  console.log('share request received! processing data');
  try {
    content = await request.json();
  } catch {
    return new Response(null, {
      status: 400,
      statusText: 'Bad Request (Invalid JSON)',
    });
  }

  //validate input
  const valid = validator.validate(content);

  if (!valid.valid) {
    console.log(valid.errors);
    return new Response(null, {status: 400, statusText: 'Bad Request'});
  }

  //save to kv
  try {
    const key = uuid();
    const data = pako.deflate(JSON.stringify(content));
    await WFPSIM_KV.put(key, data.buffer, {
      expirationTtl: 60 * 60 * 24 * 31,
    }); //7 days
    return new Response(key, {status: 200});
  } catch (err) {
    console.error(`KV returned error: ${err}`);
    return new Response('put failed', {status: 500});
  }
}
