import { Button, Group, Stack, TextInput, Table, Text, Title } from '@mantine/core';
import { useState } from 'react';
import { notifications } from '@mantine/notifications';
import { useAuth } from '../../context/AuthContext';
import {
  ClinicContact,
  loadClinicContacts,
  saveClinicContacts,
} from '../../help/clinicContacts';

export function ClinicContactsPanel() {
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';
  const [contacts, setContacts] = useState<ClinicContact[]>(() => loadClinicContacts());
  const [editing, setEditing] = useState(false);

  const updateRow = (index: number, field: keyof ClinicContact, value: string) => {
    setContacts((current) =>
      current.map((row, i) => (i === index ? { ...row, [field]: value } : row)),
    );
  };

  const handleSave = () => {
    saveClinicContacts(contacts);
    setEditing(false);
    notifications.show({
      title: 'Contacts saved',
      message: 'Stored locally on this browser at the clinic.',
      color: 'green',
    });
  };

  return (
    <Stack gap="md" mt="md">
      <Group justify="space-between">
        <Title order={4}>Your clinic contacts</Title>
        {isAdmin && !editing && (
          <Button size="xs" variant="light" onClick={() => setEditing(true)}>
            Edit contacts
          </Button>
        )}
        {isAdmin && editing && (
          <Group gap="xs">
            <Button size="xs" variant="default" onClick={() => {
              setContacts(loadClinicContacts());
              setEditing(false);
            }}>
              Cancel
            </Button>
            <Button size="xs" onClick={handleSave}>
              Save
            </Button>
          </Group>
        )}
      </Group>

      {!isAdmin && (
        <Text size="sm" c="dimmed">
          Ask an administrator to fill in emergency contacts. They are saved on the clinic server browser only.
        </Text>
      )}

      <Table withTableBorder>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Role</Table.Th>
            <Table.Th>Name</Table.Th>
            <Table.Th>Phone</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {contacts.map((contact, index) => (
            <Table.Tr key={contact.role}>
              <Table.Td>{contact.role}</Table.Td>
              <Table.Td>
                {editing ? (
                  <TextInput
                    size="xs"
                    value={contact.name}
                    placeholder="Name"
                    onChange={(event) => updateRow(index, 'name', event.currentTarget.value)}
                  />
                ) : (
                  contact.name || '—'
                )}
              </Table.Td>
              <Table.Td>
                {editing ? (
                  <TextInput
                    size="xs"
                    value={contact.phone}
                    placeholder="Phone"
                    onChange={(event) => updateRow(index, 'phone', event.currentTarget.value)}
                  />
                ) : (
                  contact.phone || '—'
                )}
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  );
}
