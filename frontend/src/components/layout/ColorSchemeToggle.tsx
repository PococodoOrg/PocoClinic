import { ActionIcon, useMantineColorScheme } from '@mantine/core';
import { IconMoon, IconSun } from '@tabler/icons-react';

interface ColorSchemeToggleProps {
  size?: number | 'sm' | 'md' | 'lg' | 'xl';
}

export function ColorSchemeToggle({ size = 'lg' }: ColorSchemeToggleProps) {
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  return (
    <ActionIcon
      variant="default"
      size={size}
      aria-label={colorScheme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
      onClick={() => toggleColorScheme()}
    >
      {colorScheme === 'dark' ? <IconSun size={20} /> : <IconMoon size={20} />}
    </ActionIcon>
  );
}
