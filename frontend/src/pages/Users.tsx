import { Container, Text, Title } from '@mantine/core';
import { UserList } from '../components/users/UserList';

export default function Users() {
  return (
    <Container size="lg">
      <Title order={2} mb="xs">
        Staff & badges
      </Title>
      <Text c="dimmed" mb="xl">
        Create staff accounts and print employee badges for badge + PIN sign-in.
      </Text>
      <UserList />
    </Container>
  );
}
