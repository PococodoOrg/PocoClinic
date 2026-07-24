export type HelpAudience = 'all' | 'staff' | 'admin';

export type HelpCategoryId =
  | 'getting-started'
  | 'daily-care'
  | 'account-security'
  | 'administration'
  | 'backup-restore'
  | 'troubleshooting';

export interface HelpCategory {
  id: HelpCategoryId;
  label: string;
  description: string;
}

export interface HelpAlert {
  color: string;
  title: string;
  body: string;
}

export interface HelpTable {
  headers: string[];
  rows: string[][];
}

export interface HelpSection {
  heading?: string;
  paragraphs?: string[];
  list?: string[];
  steps?: string[];
  table?: HelpTable;
  alert?: HelpAlert;
}

export interface HelpArticle {
  id: string;
  title: string;
  summary: string;
  category: HelpCategoryId;
  audience: HelpAudience;
  tags: string[];
  printable?: boolean;
  sections: HelpSection[];
}
