import { createContext, ReactNode, useContext, useEffect, useMemo, useState } from 'react';

interface PiTouchContextValue {
  isPiTouch: boolean;
}

const PiTouchContext = createContext<PiTouchContextValue>({ isPiTouch: false });

function detectPiTouch(): boolean {
  if (import.meta.env.VITE_PI_TOUCH === 'true') {
    return true;
  }

  const params = new URLSearchParams(window.location.search);
  if (params.get('pi') === '1') {
    localStorage.setItem('poco-pi-touch', '1');
    return true;
  }
  if (params.get('pi') === '0') {
    localStorage.removeItem('poco-pi-touch');
    return false;
  }
  if (localStorage.getItem('poco-pi-touch') === '1') {
    return true;
  }

  const coarse = window.matchMedia('(pointer: coarse)').matches;
  const small = window.matchMedia('(max-width: 900px)').matches;
  return coarse && small;
}

export function PiTouchProvider({ children }: { children: ReactNode }) {
  const [isPiTouch, setIsPiTouch] = useState(detectPiTouch);

  useEffect(() => {
    setIsPiTouch(detectPiTouch());
  }, []);

  useEffect(() => {
    document.documentElement.classList.toggle('pi-touch', isPiTouch);
    document.body.classList.toggle('pi-touch', isPiTouch);
  }, [isPiTouch]);

  const value = useMemo(() => ({ isPiTouch }), [isPiTouch]);

  return <PiTouchContext.Provider value={value}>{children}</PiTouchContext.Provider>;
}

export function usePiTouch() {
  return useContext(PiTouchContext);
}

export function useTouchUi() {
  const { isPiTouch } = usePiTouch();
  return {
    isPiTouch,
    buttonSize: isPiTouch ? ('xl' as const) : ('md' as const),
    fullWidth: isPiTouch,
    stepperOrientation: isPiTouch ? ('vertical' as const) : ('horizontal' as const),
    textSize: isPiTouch ? ('md' as const) : ('sm' as const),
  };
}
