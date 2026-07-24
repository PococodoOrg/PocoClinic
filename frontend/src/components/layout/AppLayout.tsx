import {
  ActionIcon,
  AppShell,
  Box,
  Burger,
  Button,
  Container,
  Group,
  NavLink,
  ScrollArea,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { useDisclosure, useMediaQuery } from '@mantine/hooks';
import { IconHelp } from '@tabler/icons-react';
import { ReactNode } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useHelp } from '../../context/HelpContext';
import { NAV_BREAKPOINT, TABLET_MEDIA_QUERY } from '../../layout/breakpoints';
import { HelpDrawer } from '../help/HelpDrawer';
import { PageHelpTip } from '../help/PageHelpTip';
import { getPageHelpTip } from '../../help/pageTips';
import { PageBreadcrumbs } from './PageBreadcrumbs';
import { ColorSchemeToggle } from './ColorSchemeToggle';
import { isNavActive, navSections } from './navConfig';

interface AppLayoutProps {
  children: ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const { openDrawer } = useHelp();
  const [opened, { toggle, close }] = useDisclosure();
  const isTablet = useMediaQuery(TABLET_MEDIA_QUERY);

  const handleLogout = () => {
    logout().then(() => navigate('/login'));
  };

  const visibleSections = navSections
    .map((section) => ({
      ...section,
      items: section.items.filter((item) => !item.adminOnly || user?.role === 'admin'),
    }))
    .filter((section) => section.items.length > 0);

  const pageTip = getPageHelpTip(location.pathname);

  return (
    <>
      <a href="#main-content" className="skip-link">
        Skip to main content
      </a>
      <AppShell
        header={{ height: 64 }}
        navbar={{
          width: 280,
          breakpoint: NAV_BREAKPOINT,
          collapsed: { mobile: !opened },
        }}
        padding="md"
      >
        <AppShell.Header>
          <Group h="100%" px="md" justify="space-between" wrap="nowrap">
            <Group gap="sm" wrap="nowrap">
              <Burger
                opened={opened}
                onClick={toggle}
                hiddenFrom={NAV_BREAKPOINT}
                size="md"
                aria-label={opened ? 'Close navigation' : 'Open navigation'}
              />
              <button
                type="button"
                className="app-header-brand"
                onClick={() => {
                  close();
                  navigate('/patients');
                }}
                aria-label="PocoClinic home — patients"
              >
                <Title order={4} component="span">
                  PocoClinic
                </Title>
              </button>
            </Group>
            <Group gap="xs" wrap="nowrap">
              <ActionIcon
                variant="default"
                size={isTablet ? 48 : 'lg'}
                aria-label="Open help"
                onClick={openDrawer}
                className="help-no-print"
              >
                <IconHelp size={20} />
              </ActionIcon>
              <ColorSchemeToggle size={isTablet ? 48 : undefined} />
              {user && (
                <Text size="sm" c="dimmed" visibleFrom="xs" lineClamp={1} maw={140}>
                  {user.name}
                </Text>
              )}
              <Button
                variant="default"
                size={isTablet ? 'md' : 'compact-sm'}
                onClick={handleLogout}
                className="touch-control"
              >
                Sign out
              </Button>
            </Group>
          </Group>
        </AppShell.Header>

        <AppShell.Navbar p="md" className="app-shell-nav" aria-label="Main navigation">
          <ScrollArea flex={1}>
            <Stack gap="lg">
              {visibleSections.map((section) => (
                <Box key={section.title}>
                  <Text size="xs" tt="uppercase" fw={700} c="dimmed" mb="xs">
                    {section.title}
                  </Text>
                  <Stack gap={6}>
                    {section.items.map((item) => {
                      const Icon = item.icon;
                      const active = isNavActive(location.pathname, item.path);
                      return (
                        <NavLink
                          key={item.path}
                          component={Link}
                          to={item.path}
                          label={item.label}
                          leftSection={<Icon size={20} stroke={active ? 2 : 1.5} />}
                          active={active}
                          onClick={close}
                          aria-current={active ? 'page' : undefined}
                        />
                      );
                    })}
                  </Stack>
                </Box>
              ))}
            </Stack>
          </ScrollArea>
          {user && (
            <Box pt="md" mt="auto">
              <Text size="sm" fw={600} lineClamp={1}>
                {user.name}
              </Text>
              <Text size="xs" c="dimmed" tt="capitalize">
                {user.role}
              </Text>
            </Box>
          )}
        </AppShell.Navbar>

        <AppShell.Main id="main-content" tabIndex={-1} className="app-main-content">
          <Container size="xl" px={isTablet ? 'sm' : 'md'}>
            {!isTablet && <PageBreadcrumbs />}
            {pageTip && (
              <PageHelpTip
                tipId={pageTip.tipId}
                title={pageTip.title}
                message={pageTip.message}
                articleId={pageTip.articleId}
              />
            )}
            {children}
          </Container>
        </AppShell.Main>
        <HelpDrawer />
      </AppShell>
    </>
  );
}
