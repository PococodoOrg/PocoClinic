import { Button, ButtonProps, PolymorphicComponentProps } from '@mantine/core';
import { useTouchUi } from '../context/PiTouchContext';

type TouchPrimaryButtonProps = PolymorphicComponentProps<'button', ButtonProps>;

export function TouchPrimaryButton({ size, fullWidth, ...props }: TouchPrimaryButtonProps) {
  const touch = useTouchUi();
  return (
    <Button
      size={size ?? touch.buttonSize}
      fullWidth={fullWidth ?? touch.fullWidth}
      radius="md"
      className="touch-primary-button"
      {...props}
    />
  );
}
