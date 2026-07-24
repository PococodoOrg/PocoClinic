import { Badge, Card, Group, Stack, Text, Title, Tooltip } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { SystemStatus } from '../../api/admin';
import { AdminTask, adminTasksNeedAttention, buildAdminTasks } from './adminTasks';

interface AdminTaskBoardProps {
  status: SystemStatus;
}

function taskColor(status: AdminTask['status']): string {
  switch (status) {
    case 'urgent':
      return 'red';
    case 'due':
      return 'yellow';
    default:
      return 'green';
  }
}

function taskLabel(task: AdminTask): string {
  if (task.status === 'ok') {
    return `${task.label} ✓`;
  }
  return task.label;
}

export function AdminTaskBoard({ status }: AdminTaskBoardProps) {
  const navigate = useNavigate();
  const tasks = buildAdminTasks(status);
  const needsAttention = adminTasksNeedAttention(tasks);

  const openTask = (task: AdminTask) => {
    if (!task.path) {
      return;
    }
    const [pathname, search] = task.path.split('?');
    navigate({ pathname, search: search ? `?${search}` : undefined });
  };

  return (
    <Card withBorder padding="md">
      <Stack gap="sm">
        <div>
          <Title order={4}>Admin chores</Title>
          <Text size="sm" c="dimmed">
            {needsAttention
              ? 'Tap a badge when you complete a task — green means you are caught up for now.'
              : 'All recurring chores are caught up. Nice work.'}
          </Text>
        </div>
        <Group gap="xs">
          {tasks.map((task) => (
            <Tooltip key={task.id} label={task.hint} withArrow>
              <Badge
                size="lg"
                variant={task.status === 'ok' ? 'outline' : 'filled'}
                color={taskColor(task.status)}
                style={{ cursor: task.path ? 'pointer' : 'default', textTransform: 'none' }}
                onClick={() => openTask(task)}
              >
                {taskLabel(task)}
              </Badge>
            </Tooltip>
          ))}
        </Group>
      </Stack>
    </Card>
  );
}
