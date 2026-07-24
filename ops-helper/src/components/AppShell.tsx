import { ReactNode } from 'react';
import { usePiTouch } from '../context/PiTouchContext';
import { HelperLayout } from './HelperLayout';
import { PiTouchLayout } from './PiTouchLayout';

export function AppShell({ children }: { children: ReactNode }) {
  const { isPiTouch } = usePiTouch();
  if (isPiTouch) {
    return <PiTouchLayout>{children}</PiTouchLayout>;
  }
  return <HelperLayout>{children}</HelperLayout>;
}
