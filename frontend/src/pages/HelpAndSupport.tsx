import {

  Badge,

  Button,

  Grid,

  Group,

  NavLink,

  Paper,

  ScrollArea,

  Stack,

  Text,

  TextInput,

  Title,

} from '@mantine/core';

import { IconArrowLeft, IconPrinter, IconSearch, IconStar } from '@tabler/icons-react';

import { useEffect, useMemo, useState } from 'react';

import { useNavigate, useParams, useSearchParams } from 'react-router-dom';

import { useAuth } from '../context/AuthContext';

import { useHelp } from '../context/HelpContext';

import { HelpAdminBanner } from '../components/help/HelpAdminBanner';

import { HelpArticleBody } from '../components/help/HelpArticleBody';

import { HelpArticleCard } from '../components/help/HelpArticleCard';

import { HelpRelatedArticles } from '../components/help/HelpRelatedArticles';

import {

  getArticleById,

  getCategoryById,

  getFeaturedArticleIds,

  getRelatedArticles,

  groupArticlesByCategory,

  helpCategories,

} from '../help/content';

import { filterHelpArticles } from '../help/search';

import { HelpCategoryId } from '../help/types';



export default function HelpAndSupport() {

  const { articleId: routeArticleId } = useParams<{ articleId?: string }>();

  const [searchParams, setSearchParams] = useSearchParams();

  const navigate = useNavigate();

  const { user } = useAuth();

  const { pendingArticleId, clearPendingArticle } = useHelp();



  const isAdmin = user?.role === 'admin';

  const categoryParam = searchParams.get('category') as HelpCategoryId | null;

  const queryParam = searchParams.get('q') ?? '';



  const [query, setQuery] = useState(queryParam);

  const selectedId = routeArticleId ?? pendingArticleId ?? null;



  useEffect(() => {

    if (pendingArticleId && !routeArticleId) {

      navigate(`/help/${pendingArticleId}`, { replace: true });

      clearPendingArticle();

    }

  }, [pendingArticleId, routeArticleId, navigate, clearPendingArticle]);



  useEffect(() => {

    setQuery(queryParam);

  }, [queryParam]);



  const filteredArticles = useMemo(

    () =>

      filterHelpArticles({

        query,

        categoryId: categoryParam ?? undefined,

        isAdmin,

      }),

    [query, categoryParam, isAdmin],

  );



  const groupedArticles = useMemo(

    () => groupArticlesByCategory(filteredArticles),

    [filteredArticles],

  );



  const featuredIds = useMemo(() => getFeaturedArticleIds(isAdmin), [isAdmin]);

  const adminFeaturedIds = ['admin-setup', 'admin-dashboard', 'daily-backup', 'staff-management'];



  const selectedArticle = selectedId ? getArticleById(selectedId) : undefined;

  const relatedArticles = useMemo(

    () => (selectedArticle ? getRelatedArticles(selectedArticle.id, isAdmin) : []),

    [selectedArticle, isAdmin],

  );



  useEffect(() => {

    if (selectedArticle && selectedArticle.audience === 'admin' && !isAdmin) {

      navigate('/help', { replace: true });

    }

  }, [selectedArticle, isAdmin, navigate]);



  const selectCategory = (categoryId: HelpCategoryId | null) => {

    const next = new URLSearchParams(searchParams);

    if (categoryId) {

      next.set('category', categoryId);

    } else {

      next.delete('category');

    }

    setSearchParams(next);

    navigate('/help');

  };



  const applySearch = (value: string) => {

    setQuery(value);

    const next = new URLSearchParams(searchParams);

    if (value.trim()) {

      next.set('q', value);

    } else {

      next.delete('q');

    }

    setSearchParams(next);

  };



  const openArticle = (id: string) => {

    navigate(`/help/${id}`);

  };



  const handlePrint = () => {

    window.print();

  };



  const showBrowseHome = !selectedArticle && !categoryParam && !query.trim();



  return (

    <Stack gap="lg" className="help-hub">

      <Group justify="space-between" align="flex-start" className="help-no-print">

        <div>

          <Title order={2}>Help center</Title>

          <Text c="dimmed" mt="xs" maw={600}>

            Guides for daily clinic work, administration, backups, and troubleshooting — on your private LAN, no internet required.

          </Text>

        </div>

        {selectedArticle?.printable && (

          <Button

            variant="light"

            leftSection={<IconPrinter size={16} />}

            onClick={handlePrint}

          >

            Print

          </Button>

        )}

      </Group>



      {isAdmin && showBrowseHome && (

        <HelpAdminBanner featuredIds={adminFeaturedIds} />

      )}



      <Grid gutter="lg">

        <Grid.Col span={{ base: 12, md: 4 }}>

          <Paper withBorder p="md" className="help-sidebar help-no-print">

            <Stack gap="md">

              <TextInput

                placeholder="Search help articles…"

                leftSection={<IconSearch size={16} />}

                value={query}

                onChange={(event) => applySearch(event.currentTarget.value)}

                aria-label="Search help articles"

              />



              <Stack gap={4}>

                <NavLink

                  label="All topics"

                  active={!categoryParam && !query.trim()}

                  onClick={() => {

                    applySearch('');

                    selectCategory(null);

                  }}

                />

                {helpCategories.map((category) => (

                  <NavLink

                    key={category.id}

                    label={category.label}

                    description={category.description}

                    active={categoryParam === category.id}

                    onClick={() => selectCategory(category.id)}

                  />

                ))}

              </Stack>



              {showBrowseHome && (

                <Stack gap="xs">

                  <Group gap={6}>

                    <IconStar size={14} />

                    <Text size="xs" tt="uppercase" fw={700} c="dimmed">Start here</Text>

                  </Group>

                  {featuredIds.slice(0, 3).map((id) => {

                    const article = getArticleById(id);

                    if (!article) {

                      return null;

                    }

                    return (

                      <HelpArticleCard

                        key={id}

                        article={article}

                        compact

                        active={selectedId === id}

                        onClick={() => openArticle(id)}

                      />

                    );

                  })}

                </Stack>

              )}

            </Stack>

          </Paper>

        </Grid.Col>



        <Grid.Col span={{ base: 12, md: 8 }}>

          {selectedArticle ? (

            <Paper withBorder p="xl" className="help-article-view">

              <Stack gap="md">

                <Group className="help-no-print">

                  <Button

                    variant="subtle"

                    leftSection={<IconArrowLeft size={16} />}

                    onClick={() => navigate('/help')}

                  >

                    All articles

                  </Button>

                </Group>



                <Group gap="xs">

                  <Badge variant="light">

                    {getCategoryById(selectedArticle.category)?.label}

                  </Badge>

                  {selectedArticle.audience === 'admin' && (

                    <Badge color="grape" variant="light">

                      Administrators

                    </Badge>

                  )}

                  {selectedArticle.printable && (

                    <Badge color="gray" variant="outline" className="help-no-print">

                      Printable

                    </Badge>

                  )}

                </Group>



                <Title order={2}>{selectedArticle.title}</Title>

                <Text c="dimmed">{selectedArticle.summary}</Text>



                <HelpArticleBody

                  sections={selectedArticle.sections}

                  showContacts={selectedArticle.id === 'clinic-contacts'}

                />



                <HelpRelatedArticles

                  articles={relatedArticles}

                  onSelect={openArticle}

                />

              </Stack>

            </Paper>

          ) : (

            <Stack gap="lg">

              <Text fw={600}>

                {categoryParam

                  ? getCategoryById(categoryParam)?.label

                  : query.trim()

                    ? `Results for “${query.trim()}”`

                    : 'Browse help topics'}

              </Text>



              {categoryParam && (

                <Text size="sm" c="dimmed">

                  {getCategoryById(categoryParam)?.description}

                </Text>

              )}



              {filteredArticles.length === 0 ? (

                <Text c="dimmed">No articles match your search. Try another term or category.</Text>

              ) : query.trim() || categoryParam ? (

                <ScrollArea.Autosize mah={700}>

                  <Stack gap="sm">

                    {filteredArticles.map((article) => (

                      <HelpArticleCard

                        key={article.id}

                        article={article}

                        onClick={() => openArticle(article.id)}

                      />

                    ))}

                  </Stack>

                </ScrollArea.Autosize>

              ) : (

                groupedArticles.map(({ category, articles }) => (

                  <Stack key={category.id} gap="sm">

                    <div>

                      <Text fw={600}>{category.label}</Text>

                      <Text size="sm" c="dimmed">{category.description}</Text>

                    </div>

                    <Stack gap="xs">

                      {articles.map((article) => (

                        <HelpArticleCard

                          key={article.id}

                          article={article}

                          onClick={() => openArticle(article.id)}

                        />

                      ))}

                    </Stack>

                  </Stack>

                ))

              )}

            </Stack>

          )}

        </Grid.Col>

      </Grid>

    </Stack>

  );

}


