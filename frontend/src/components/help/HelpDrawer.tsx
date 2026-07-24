import {

  Button,

  Divider,

  Drawer,

  ScrollArea,

  Stack,

  Text,

  TextInput,

  Title,

} from '@mantine/core';

import { IconBook, IconSearch, IconSettings } from '@tabler/icons-react';

import { useMemo, useState } from 'react';

import { useLocation, useNavigate } from 'react-router-dom';

import { useAuth } from '../../context/AuthContext';

import { useHelp } from '../../context/HelpContext';

import { getArticleById } from '../../help/content';

import { getContextualArticleIds } from '../../help/contextual';

import { searchHelpArticles } from '../../help/search';

import { HelpArticleCard } from './HelpArticleCard';



export function HelpDrawer() {

  const { drawerOpen, closeDrawer } = useHelp();

  const { user } = useAuth();

  const location = useLocation();

  const navigate = useNavigate();

  const [query, setQuery] = useState('');



  const isAdmin = user?.role === 'admin';

  const onAdminPage = location.pathname.startsWith('/admin');



  const contextualArticles = useMemo(() => {

    return getContextualArticleIds(location.pathname)

      .map((id) => getArticleById(id))

      .filter((article) => article && (article.audience === 'all' || isAdmin));

  }, [location.pathname, isAdmin]);



  const searchResults = useMemo(() => searchHelpArticles(query, isAdmin), [query, isAdmin]);



  const openArticle = (articleId: string) => {

    closeDrawer();

    setQuery('');

    navigate(`/help/${articleId}`);

  };



  const goToHub = (category?: string) => {

    closeDrawer();

    setQuery('');

    navigate(category ? `/help?category=${category}` : '/help');

  };



  return (

    <Drawer

      opened={drawerOpen}

      onClose={() => {

        closeDrawer();

        setQuery('');

      }}

      title="Help"

      position="right"

      size="md"

    >

      <Stack gap="md">

        <TextInput

          placeholder="Search help…"

          leftSection={<IconSearch size={16} />}

          value={query}

          onChange={(event) => setQuery(event.currentTarget.value)}

          aria-label="Search help articles"

        />



        {query.trim() ? (

          <ScrollArea.Autosize mah={400}>

            <Stack gap="xs">

              {searchResults.length === 0 ? (

                <Text size="sm" c="dimmed">No matching articles. Try different words or open full Help.</Text>

              ) : (

                searchResults.map((article) => (

                  <HelpArticleCard

                    key={article.id}

                    article={article}

                    compact

                    onClick={() => openArticle(article.id)}

                  />

                ))

              )}

            </Stack>

          </ScrollArea.Autosize>

        ) : (

          <>

            {isAdmin && onAdminPage && (

              <Stack gap="xs">

                <Title order={5}>Administrator</Title>

                <Button

                  variant="light"

                  leftSection={<IconSettings size={16} />}

                  onClick={() => {

                    closeDrawer();

                    navigate('/admin');

                  }}

                >

                  Admin dashboard & guide

                </Button>

                <Button variant="subtle" onClick={() => openArticle('admin-setup')}>

                  First-time setup

                </Button>

              </Stack>

            )}



            {contextualArticles.length > 0 && (

              <Stack gap="xs">

                <Title order={5}>Help for this page</Title>

                {contextualArticles.map((article) => (

                  article && (

                    <HelpArticleCard

                      key={article.id}

                      article={article}

                      compact

                      onClick={() => openArticle(article.id)}

                    />

                  )

                ))}

              </Stack>

            )}



            <Divider />



            <Stack gap="xs">

              <Title order={5}>Quick links</Title>

              <Button variant="light" onClick={() => openArticle('sign-in')}>

                Sign in with badge & PIN

              </Button>

              <Button variant="light" onClick={() => openArticle('patient-list')}>

                Find patients

              </Button>

              {isAdmin && (

                <>

                  <Button variant="light" onClick={() => openArticle('daily-backup')}>

                    Daily USB backup

                  </Button>

                  <Button variant="light" leftSection={<IconBook size={16} />} onClick={() => goToHub('administration')}>

                    Administration articles

                  </Button>

                </>

              )}

              <Button variant="default" onClick={() => goToHub()}>

                Open full Help center

              </Button>

            </Stack>

          </>

        )}

      </Stack>

    </Drawer>

  );

}


