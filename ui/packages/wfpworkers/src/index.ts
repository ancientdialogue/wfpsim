import {Router} from 'itty-router';
import {handleAssets} from './assets';
import {handleInjectHead, handlePreview} from './preview';
import {handleShare, handleView} from './share';
import {handleWasm} from './wasm';

const router = Router();

router.post('/api/share', handleShare);
router.get('/api/share/:key', handleView);
router.get('/api/share/db/:key', handleView);
router.get('/api/share/sh/:key', handleView);
router.get('/api/assets/*', handleAssets);
router.get('/api/wasm/*', handleWasm);

// rewrite doc head
router.get('/sh/:key', handleInjectHead);
router.get('/api/preview/:key', handlePreview);

addEventListener('fetch', (event) => {
  event.respondWith(router.handle(event.request, event));
});
