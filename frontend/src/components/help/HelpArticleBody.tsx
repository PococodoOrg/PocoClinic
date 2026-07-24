import {
  Alert,
  List,
  Stack,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { HelpSection } from '../../help/types';
import { ClinicContactsPanel } from './ClinicContactsPanel';

interface HelpArticleBodyProps {
  sections: HelpSection[];
  showContacts?: boolean;
}

export function HelpArticleBody({ sections, showContacts }: HelpArticleBodyProps) {
  return (
    <Stack gap="lg" className="help-article-body">
      {sections.map((section, index) => (
        <Stack key={index} gap="sm">
          {section.heading && <Title order={4}>{section.heading}</Title>}
          {section.paragraphs?.map((paragraph) => (
            <Text key={paragraph} size="sm">
              {paragraph}
            </Text>
          ))}
          {section.list && (
            <List size="sm" spacing="xs">
              {section.list.map((item) => (
                <List.Item key={item}>{item}</List.Item>
              ))}
            </List>
          )}
          {section.steps && (
            <List type="ordered" size="sm" spacing="xs">
              {section.steps.map((item) => (
                <List.Item key={item}>{item}</List.Item>
              ))}
            </List>
          )}
          {section.table && (
            <Table withTableBorder striped highlightOnHover>
              <Table.Thead>
                <Table.Tr>
                  {section.table.headers.map((header) => (
                    <Table.Th key={header}>{header}</Table.Th>
                  ))}
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {section.table.rows.map((row, rowIndex) => (
                  <Table.Tr key={rowIndex}>
                    {row.map((cell) => (
                      <Table.Td key={cell}>{cell}</Table.Td>
                    ))}
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          )}
          {section.alert && (
            <Alert color={section.alert.color} title={section.alert.title}>
              {section.alert.body}
            </Alert>
          )}
        </Stack>
      ))}
      {showContacts && <ClinicContactsPanel />}
    </Stack>
  );
}
