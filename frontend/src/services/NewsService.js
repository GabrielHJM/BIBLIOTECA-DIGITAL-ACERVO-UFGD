import axios from 'axios';

const GNEWS_API_KEY = process.env.VUE_APP_GNEWS_API_KEY || '';
const GNEWS_BASE_URL = 'https://gnews.io/api/v4/search';
const CACHE_TTL_MS = 10 * 60 * 1000; // 10 minutes in milliseconds

/**
 * Servico para buscar noticias na API do GNews.
 * Implementa cache local de 10 minutos para economizar limite de 100 requisições/dia.
 */
export default {
	async buscarNoticias(query = 'ciência OU tecnologia OU educação', max = 5, lang = 'pt') {
		if (!GNEWS_API_KEY) {
			console.warn('GNEWS_API_KEY não está configurada no .env.local. Retornando dados mockados para desenvolvimento.');
			return this.getMockNews();
		}

		const cacheKey = `gnews_cache_v2_${query}_${lang}_${max}`;
		const cachedDataStr = localStorage.getItem(cacheKey);

		if (cachedDataStr) {
			try {
				const cache = JSON.parse(cachedDataStr);
				const now = new Date().getTime();
				if (now - cache.timestamp < CACHE_TTL_MS) {
					console.log('Retornando notícias do cache local (GNews).');
					return cache.data;
				}
			} catch (e) {
				console.error('Erro ao ler cache do GNews', e);
			}
		}

		try {
			const url = `${GNEWS_BASE_URL}?q=${encodeURIComponent(query)}&lang=${lang}&max=${max}&apikey=${GNEWS_API_KEY}`;
			const response = await axios.get(url);

			const articles = response.data.articles || [];

			// Mapeamento extra caso o JSON venha traduzido (ex: Chrome autotranslate) ou padrão Gnews
			const normalizedArticles = articles.map(a => ({
				title: a.title || a.título,
				description: a.description || a.descrição,
				url: a.url,
				image: (a.image !== 'None' ? a.image : null) || (a.imagem !== 'None' ? a.imagem : null),
				source: { name: (a.source && a.source.name) || (a.fonte && a.fonte.nome) || 'GNews' }
			}));

			// Salva no cache com limite de tempo
			const cacheObj = {
				timestamp: new Date().getTime(),
				data: normalizedArticles
			};
			localStorage.setItem(cacheKey, JSON.stringify(cacheObj));

			return normalizedArticles;
		} catch (error) {
			console.error('Erro ao buscar notícias no GNews API:', error);
			// Se der erro (ex: limite atingido), tentamos puxar do cache mesmo se estiver expirado, ou dados mockados
			if (cachedDataStr) {
				return JSON.parse(cachedDataStr).data;
			}
			return this.getMockNews();
		}
	},

	getMockNews() {
		return [
			{
				title: "Nova descoberta na computação quântica acelera processamento de IA.",
				description: "Pesquisadores conseguiram manter qubits estáveis por tempo recorde em temperatura ambiente.",
				url: "https://gnews.io/",
				image: "https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=400",
				source: { name: "Tech Science News" }
			},
			{
				title: "MEC anuncia ampliação de bibliotecas digitais no ensino médio.",
				description: "O programa visa levar acervos virtuais para mais de 5.000 escolas públicas até o fim do semestre.",
				url: "https://gnews.io/",
				image: "https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=400",
				source: { name: "Educação Hoje" }
			},
			{
				title: "Como a Inteligência Artificial está transformando a medicina diagnóstica.",
				description: "Novos algoritmos conseguem detectar patologias em exames de imagem com 99% de precisão.",
				url: "https://gnews.io/",
				image: "https://images.unsplash.com/photo-1532187863486-abf9dbad1b69?w=400",
				source: { name: "Health Tech Review" }
			},
			{
				title: "Telescópio James Webb encontra possível exoplaneta habitável.",
				description: "A assinatura d'água foi confirmada na atmosfera de um planeta parecido com a Terra.",
				url: "https://gnews.io/",
				image: "https://images.unsplash.com/photo-1614730321146-b6fa6a46bcb4?w=400",
				source: { name: "Space Observer" }
			},
			{
				title: "Evolução do 6G: Primeiros testes começam em laboratórios globais.",
				description: "Empresas começam a padronizar a próxima geração de conectividade móvel.",
				url: "https://gnews.io/",
				image: "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=400",
				source: { name: "Connectivity Weekly" }
			}
		];
	}
};
