import { Button, Card, List, Stack, Text, Title } from '@mantine/core';
import { Link } from 'react-router-dom';

export function ChecklistPage() {
  return (
    <Stack gap="lg">
      <Card withBorder padding="xl" radius="lg">
        <Title order={2} mb="xs">
          Daily backup checklist
        </Title>
        <Text c="dimmed" mb="lg" className="no-print">
          Print this page and keep it in the clinic binder next to your USB drives.
        </Text>

        <List spacing="md" size="lg" type="ordered">
          <List.Item>PocoClinic is running and staff can sign in.</List.Item>
          <List.Item>Insert today&apos;s labeled USB drive into this computer.</List.Item>
          <List.Item>
            Open the Backup Helper on this computer: <strong>http://127.0.0.1:9090</strong>
          </List.Item>
          <List.Item>Follow <strong>Daily backup</strong> and click <strong>Create backup now</strong>.</List.Item>
          <List.Item>Copy the backup file onto the USB drive.</List.Item>
          <List.Item>Safely remove USB and store in locked location.</List.Item>
          <List.Item>
            Initials / date: _______________________________
          </List.Item>
        </List>

        <Stack gap="sm" mt="xl" className="no-print">
          <Button onClick={() => window.print()}>Print this checklist</Button>
          <Button component={Link} to="/backup" variant="light">
            Start backup guide
          </Button>
        </Stack>
      </Card>
    </Stack>
  );
}
