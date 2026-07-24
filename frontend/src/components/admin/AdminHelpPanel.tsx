import {
  Accordion,
  Badge,
  Button,
  Card,
  Checkbox,
  Group,
  Paper,
  Progress,
  SimpleGrid,
  Stack,
  Tabs,
  Text,
  Title,
} from '@mantine/core';
import {
  IconBook,
  IconCheck,
  IconExternalLink,
  IconListCheck,
  IconRotateClockwise,
} from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminGuideSections, setupChecklistTasks } from '../../help/adminGuide';
import { loadSetupChecklist, saveSetupChecklist, toggleSetupStep } from '../../help/adminSetupChecklist';

function AdminGuideTaskRow({
  title,
  description,
  route,
  articleId,
  onNavigate,
}: {
  title: string;
  description: string;
  route?: string;
  articleId?: string;
  onNavigate: (path: string) => void;
}) {
  return (
    <Paper withBorder p="sm" radius="md" className="admin-guide-task">
      <Text size="sm" fw={600}>{title}</Text>
      <Text size="xs" c="dimmed" mt={4}>{description}</Text>
      <Group gap="xs" mt="xs">
        {route && (
          <Button size="compact-xs" variant="light" onClick={() => onNavigate(route)}>
            Open
          </Button>
        )}
        {articleId && (
          <Button
            size="compact-xs"
            variant="default"
            leftSection={<IconExternalLink size={12} />}
            onClick={() => onNavigate(`/help/${articleId}`)}
          >
            Guide
          </Button>
        )}
      </Group>
    </Paper>
  );
}

export function AdminHelpPanel() {
  const navigate = useNavigate();
  const setupTasks = setupChecklistTasks();
  const referenceSections = adminGuideSections.filter((section) => section.id !== 'initial-setup');
  const [checklist, setChecklist] = useState<Record<string, boolean>>(() => loadSetupChecklist());

  const completedCount = useMemo(
    () => setupTasks.filter((task) => checklist[task.id]).length,
    [checklist, setupTasks],
  );
  const progress = setupTasks.length === 0 ? 0 : Math.round((completedCount / setupTasks.length) * 100);
  const setupComplete = progress === 100;

  const handleToggle = (stepId: string, checked: boolean) => {
    setChecklist(toggleSetupStep(stepId, checked));
  };

  const resetChecklist = () => {
    saveSetupChecklist({});
    setChecklist({});
  };

  return (
    <Card withBorder padding="lg" className="admin-help-panel">
      <Group justify="space-between" mb="lg" align="flex-start" wrap="wrap">
        <Stack gap={4}>
          <Title order={3}>Administrator guide</Title>
          <Text size="sm" c="dimmed" maw={640}>
            First-time setup and every ongoing admin task — with links into the app and detailed help articles.
          </Text>
        </Stack>
        <Group>
          <Button
            variant="light"
            leftSection={<IconBook size={16} />}
            onClick={() => navigate('/help?category=administration')}
          >
            All admin articles
          </Button>
          <Button variant="default" onClick={() => navigate('/help/admin-setup')}>
            Full setup guide
          </Button>
        </Group>
      </Group>

      <Tabs defaultValue={setupComplete ? 'tasks' : 'setup'}>
        <Tabs.List mb="md">
          <Tabs.Tab value="setup" leftSection={<IconListCheck size={16} />}>
            Setup checklist
            {!setupComplete && (
              <Badge size="xs" variant="light" ml={6}>{setupTasks.length - completedCount} left</Badge>
            )}
          </Tabs.Tab>
          <Tabs.Tab value="tasks" leftSection={<IconBook size={16} />}>
            Task library
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="setup">
          <Stack gap="md">
            <Group justify="space-between">
              <Group gap="sm">
                <Text fw={600} size="sm">First-time clinic setup</Text>
                <Badge variant="light" color={setupComplete ? 'green' : 'blue'}>
                  {completedCount} / {setupTasks.length}
                </Badge>
              </Group>
              {completedCount > 0 && (
                <Button
                  size="xs"
                  variant="subtle"
                  leftSection={<IconRotateClockwise size={14} />}
                  onClick={resetChecklist}
                >
                  Reset checklist
                </Button>
              )}
            </Group>

            <Progress value={progress} size="lg" radius="xl" />

            {setupComplete && (
              <Paper withBorder p="md" bg="var(--mantine-color-green-light)">
                <Group gap="xs">
                  <IconCheck size={18} />
                  <Text size="sm" fw={600}>Setup checklist complete</Text>
                </Group>
                <Text size="sm" c="dimmed" mt={4}>
                  Continue with daily operations in the Task library tab and Maintenance reminders below.
                </Text>
              </Paper>
            )}

            <Stack gap="sm">
              {setupTasks.map((task, index) => {
                const done = Boolean(checklist[task.id]);
                return (
                  <Paper
                    key={task.id}
                    withBorder
                    p="md"
                    radius="md"
                    className={done ? 'admin-setup-step admin-setup-step--done' : 'admin-setup-step'}
                  >
                    <Group wrap="nowrap" align="flex-start">
                      <Checkbox
                        checked={done}
                        onChange={(event) => handleToggle(task.id, event.currentTarget.checked)}
                        aria-label={task.title}
                        mt={2}
                      />
                      <Stack gap={4} style={{ flex: 1 }}>
                        <Group gap="xs">
                          <Badge size="sm" variant="outline" circle>
                            {index + 1}
                          </Badge>
                          <Text size="sm" fw={600} td={done ? 'line-through' : undefined} c={done ? 'dimmed' : undefined}>
                            {task.title}
                          </Text>
                        </Group>
                        <Text size="xs" c="dimmed">{task.description}</Text>
                        <Group gap="xs">
                          {task.route && (
                            <Button size="compact-xs" variant="light" onClick={() => navigate(task.route!)}>
                              Open
                            </Button>
                          )}
                          {task.articleId && (
                            <Button
                              size="compact-xs"
                              variant="subtle"
                              onClick={() => navigate(`/help/${task.articleId}`)}
                            >
                              Read guide
                            </Button>
                          )}
                        </Group>
                      </Stack>
                    </Group>
                  </Paper>
                );
              })}
            </Stack>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="tasks">
          <Accordion variant="separated" defaultValue={referenceSections[0]?.id}>
            {referenceSections.map((section) => (
              <Accordion.Item key={section.id} value={section.id}>
                <Accordion.Control>
                  <Text fw={600}>{section.title}</Text>
                  <Text size="xs" c="dimmed">{section.description}</Text>
                </Accordion.Control>
                <Accordion.Panel>
                  <SimpleGrid cols={{ base: 1, md: 2 }} spacing="sm">
                    {section.tasks.map((task) => (
                      <AdminGuideTaskRow
                        key={task.id}
                        title={task.title}
                        description={task.description}
                        route={task.route}
                        articleId={task.articleId}
                        onNavigate={navigate}
                      />
                    ))}
                  </SimpleGrid>
                </Accordion.Panel>
              </Accordion.Item>
            ))}
          </Accordion>
        </Tabs.Panel>
      </Tabs>

      <Text size="xs" c="dimmed" mt="lg">
        Checklist progress is saved in this browser only. Printable backup checklists:{' '}
        <Text
          span
          size="xs"
          c="blue"
          style={{ cursor: 'pointer' }}
          onClick={() => navigate('/help/daily-backup')}
        >
          Daily USB backup
        </Text>
      </Text>
    </Card>
  );
}
