import {uuid} from '@cfworker/uuid';
import {IRequest} from 'itty-router';
import pako from 'pako';
import {validator} from './validation';

function getCharNames(data) {
  if (!data) {
    return '';
  }
  const names = (data.character_details ?? [])
    .map((character) => character?.name.slice(0, 12))
    .filter((name) => typeof name === 'string' && name.trim().length > 0)
    .slice(0, 4);
  console.log(names);
  const sortedNames = names.toSorted((a, b) => a.localeCompare(b));
  return sortedNames.length > 0 ? sortedNames.join('-') + '-' : '';
}

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
    const charNamePrefix = getCharNames(content);
    console.log(charNamePrefix);
    const key = charNamePrefix + uuid();
    const data = pako.deflate(JSON.stringify(content));
    await WFPSIM_KV.put(key, data.buffer, {
      expirationTtl: 60 * 60 * 24 * 180,
    }); //180 days
    return new Response(key, {status: 200});
  } catch (err) {
    console.error(`KV returned error: ${err}`);
    return new Response('put failed', {status: 500});
  }
}
