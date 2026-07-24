/**

 * Returns a safe same-origin relative path for post-login redirects.

 */

export function safeInternalPath(path: unknown, fallback = '/patients'): string {

  if (typeof path !== 'string') {

    return fallback;

  }



  let decoded = path;

  try {

    decoded = decodeURIComponent(path);

  } catch {

    return fallback;

  }



  if (

    !decoded.startsWith('/') ||

    decoded.startsWith('//') ||

    decoded.includes('://') ||

    decoded.includes('\\')

  ) {

    return fallback;

  }



  return decoded;

}

