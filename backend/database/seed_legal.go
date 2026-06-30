package database

import (
	"strings"

	"upcycleconnect/backend/models"

	"gorm.io/gorm"
)

const legalMentionsContent = `<h2>Mentions légales</h2>
<p>Conformément aux dispositions des articles 6-III et 19 de la loi n° 2004-575 du 21 juin 2004 pour la confiance dans l'économie numérique (LCEN), les présentes mentions légales sont portées à la connaissance des utilisateurs du site UpcycleConnect.</p>

<h3>1. Nature du projet</h3>
<p><strong>UpcycleConnect</strong> est un projet étudiant à but pédagogique, réalisé dans le cadre du Projet Annuel de 2e année à l'<strong>ESGI (École Supérieure de Génie Informatique)</strong>. Il ne constitue pas une activité commerciale et n'est rattaché à aucune société immatriculée. À ce titre, il n'est pas inscrit au Registre du Commerce et des Sociétés et ne dispose ni de numéro SIRET, ni de capital social.</p>
<p>Les fonctionnalités de paiement, de réservation, d'abonnement et de mise en relation présentées sur le site le sont à des fins de démonstration. Aucune transaction commerciale réelle n'est réalisée : les paiements éventuels s'effectuent dans un environnement de test.</p>

<h3>2. Éditeur du site</h3>
<ul>
<li>Nom du site : UpcycleConnect</li>
<li>Adresse : https://upcycleconnect.xyz</li>
<li>Directeur de la publication : Arthur KETCHEIAN</li>
<li>Contact : contact@upcycleconnect.xyz</li>
<li>Cadre : Projet Annuel — ESGI</li>
</ul>

<h3>3. Hébergeur</h3>
<p>Le site est hébergé par :</p>
<ul>
<li><strong>OVH SAS</strong></li>
<li>Siège social : 2 rue Kellermann, 59100 Roubaix, France</li>
<li>RCS Lille Métropole 424 761 419 00045</li>
<li>Téléphone : +33 9 72 10 10 07</li>
<li>Site : https://www.ovhcloud.com</li>
</ul>

<h3>4. Nom de domaine</h3>
<p>Le nom de domaine upcycleconnect.xyz est enregistré auprès du bureau d'enregistrement OVH SAS, dont les coordonnées figurent ci-dessus.</p>

<h3>5. Propriété intellectuelle</h3>
<p>L'ensemble des éléments du site (textes, logo, charte graphique, code source, illustrations) est protégé par le droit de la propriété intellectuelle. Toute reproduction ou représentation, totale ou partielle, sans autorisation préalable est interdite. Les contenus publiés par les utilisateurs (projets, photos, messages) restent la propriété de leurs auteurs, qui en concèdent un droit d'usage pour les seuls besoins du fonctionnement du site.</p>

<h3>6. Responsabilité</h3>
<p>S'agissant d'un projet pédagogique, le site est fourni « en l'état », sans garantie de disponibilité ni d'absence d'erreurs. L'éditeur ne saurait être tenu responsable des dommages directs ou indirects liés à l'utilisation du site.</p>

<h3>7. Données personnelles</h3>
<p>Le traitement des données personnelles est décrit dans la <a href="/confidentialite">Politique de confidentialité</a>. Pour toute question, vous pouvez écrire à dpo@upcycleconnect.xyz.</p>

<h3>8. Droit applicable</h3>
<p>Les présentes mentions sont soumises au droit français. En cas de litige, et à défaut de résolution amiable, les tribunaux français sont seuls compétents.</p>`

const legalCGUContent = `<h2>Conditions Générales d'Utilisation et de Vente</h2>
<p>Version en vigueur : juin 2026.</p>

<h3>Article 1 — Objet</h3>
<p>Les présentes conditions générales (ci-après les « CGU/CGV ») régissent l'accès et l'utilisation de la plateforme UpcycleConnect (ci-après la « Plateforme ») ainsi que les modalités des services proposés. Toute utilisation de la Plateforme implique l'acceptation pleine et entière des présentes.</p>

<h3>Article 2 — Définitions</h3>
<ul>
<li><strong>Plateforme</strong> : le site UpcycleConnect et l'ensemble de ses services.</li>
<li><strong>Utilisateur</strong> : toute personne disposant d'un compte.</li>
<li><strong>Client</strong> : utilisateur recourant à une prestation ou à un service.</li>
<li><strong>Prestataire</strong> : utilisateur professionnel proposant des prestations d'upcycling.</li>
<li><strong>Prestation</strong> : service de transformation, de réparation ou de réemploi proposé par un Prestataire.</li>
</ul>

<h3>Article 3 — Cadre pédagogique</h3>
<p>UpcycleConnect est un projet étudiant non commercial (voir <a href="/mentions-legales">Mentions légales</a>). Les services payants sont présentés à titre démonstratif et les paiements s'effectuent en environnement de test ; aucune somme réelle n'est encaissée.</p>

<h3>Article 4 — Inscription et compte</h3>
<p>La création d'un compte nécessite la fourniture d'informations exactes et la validation de l'adresse e-mail. La connexion peut également s'effectuer via Google ou Apple. L'Utilisateur est responsable de la confidentialité de ses identifiants et de toute activité réalisée depuis son compte. La Plateforme est réservée aux personnes âgées d'au moins 15 ans.</p>

<h3>Article 5 — Description des services</h3>
<p>La Plateforme propose notamment :</p>
<ul>
<li>la mise en relation entre Clients et Prestataires d'upcycling ;</li>
<li>la réservation et le paiement de prestations ;</li>
<li>la participation à des événements et formations, gratuits ou payants ;</li>
<li>des abonnements « Pro » et des campagnes de mise en avant pour les Prestataires ;</li>
<li>le dépôt d'objets et la localisation de points de collecte ;</li>
<li>un espace communautaire (forum, conseils, projets d'upcycling) ;</li>
<li>un service d'assistance par tickets.</li>
</ul>

<h3>Article 6 — Conditions financières</h3>
<p>Les prix sont indiqués en euros, toutes taxes comprises (TVA applicable de 20 %). Les paiements sont traités par le prestataire de paiement <strong>Stripe</strong>. Une commission de <strong>10 %</strong> est prélevée par la Plateforme sur le montant des prestations réglées, le solde étant reversé au Prestataire. Les abonnements Pro et les campagnes publicitaires font l'objet d'une tarification spécifique présentée avant toute souscription. Ces flux étant simulés dans le cadre du projet, aucune facturation réelle n'intervient.</p>

<h3>Article 7 — Réservation et facturation</h3>
<p>La réservation d'une prestation peut donner lieu à l'émission d'un devis, puis d'une facture, accessibles depuis l'espace personnel. La prestation est confirmée après acceptation et paiement.</p>

<h3>Article 8 — Droit de rétractation</h3>
<p>Conformément aux articles L221-18 et suivants du Code de la consommation, le Client consommateur dispose d'un délai de <strong>quatorze (14) jours</strong> à compter de la conclusion du contrat pour exercer son droit de rétractation, sans avoir à motiver sa décision ni à supporter de frais.</p>
<p>Pour l'exercer, le Client informe la Plateforme de sa décision au moyen d'une déclaration dénuée d'ambiguïté, adressée à contact@upcycleconnect.xyz.</p>
<p>Le droit de rétractation ne peut toutefois plus être exercé pour les prestations de services pleinement exécutées avant la fin du délai, lorsque l'exécution a commencé avec l'accord exprès du Client et que celui-ci a reconnu perdre son droit de rétractation une fois la prestation pleinement exécutée.</p>

<h3>Article 9 — Obligations des utilisateurs</h3>
<p>L'Utilisateur s'engage à utiliser la Plateforme de bonne foi, à ne pas publier de contenu illicite, diffamatoire, trompeur ou portant atteinte aux droits de tiers, et à respecter les autres utilisateurs, notamment au sein des espaces communautaires. Tout manquement peut entraîner la suppression des contenus concernés ou la suspension du compte.</p>

<h3>Article 10 — Rôle de la Plateforme</h3>
<p>La Plateforme agit en qualité d'intermédiaire technique de mise en relation et d'hébergeur des contenus publiés par les utilisateurs. Elle n'est pas partie aux contrats conclus entre Clients et Prestataires et ne saurait être tenue responsable de leur bonne exécution.</p>

<h3>Article 11 — Propriété intellectuelle</h3>
<p>Les éléments de la Plateforme sont protégés. Les utilisateurs conservent les droits sur les contenus qu'ils publient et concèdent à la Plateforme une licence d'utilisation pour les besoins du service.</p>

<h3>Article 12 — Données personnelles</h3>
<p>Les traitements de données à caractère personnel sont décrits dans la <a href="/confidentialite">Politique de confidentialité</a>.</p>

<h3>Article 13 — Suspension et résiliation</h3>
<p>L'Utilisateur peut supprimer son compte à tout moment. La Plateforme peut suspendre ou clôturer un compte en cas de non-respect des présentes conditions.</p>

<h3>Article 14 — Responsabilité</h3>
<p>Le service étant fourni dans un cadre pédagogique, il est proposé « en l'état », sans garantie de continuité ni d'absence d'erreurs. La responsabilité de l'éditeur ne peut être engagée au titre des dommages indirects résultant de l'utilisation de la Plateforme.</p>

<h3>Article 15 — Modification des conditions</h3>
<p>La Plateforme se réserve le droit de modifier les présentes CGU/CGV. La version applicable est celle en vigueur à la date d'utilisation de la Plateforme.</p>

<h3>Article 16 — Droit applicable et litiges</h3>
<p>Les présentes sont soumises au droit français. En cas de litige, le Client consommateur peut recourir gratuitement à un médiateur de la consommation ou à la plateforme européenne de Règlement en Ligne des Litiges (https://ec.europa.eu/consumers/odr). À défaut de résolution amiable, les tribunaux français sont compétents.</p>`

const legalPrivacyContent = `<h2>Politique de confidentialité</h2>
<p>La présente politique décrit la manière dont UpcycleConnect collecte et traite les données à caractère personnel des utilisateurs, conformément au Règlement (UE) 2016/679 (RGPD) et à la loi « Informatique et Libertés ».</p>

<h3>1. Responsable du traitement</h3>
<p>Le responsable du traitement est le projet UpcycleConnect (projet étudiant ESGI), représenté par son directeur de la publication, Arthur KETCHEIAN. Pour toute question relative à vos données : <strong>dpo@upcycleconnect.xyz</strong>.</p>

<h3>2. Données collectées</h3>
<ul>
<li>Données d'identité : nom, prénom et, le cas échéant, adresse postale et téléphone ;</li>
<li>Données de compte : adresse e-mail, mot de passe (stocké de façon chiffrée), photo de profil ;</li>
<li>Données de connexion via Google ou Apple (identifiant du compte) ;</li>
<li>Données professionnelles (pour les Prestataires) : informations d'entreprise, SIRET, profil public ;</li>
<li>Données de paiement : traitées par Stripe ; UpcycleConnect ne conserve aucune coordonnée bancaire ;</li>
<li>Contenus publiés : messages, tickets d'assistance, publications du forum, projets, photos d'objets déposés ;</li>
<li>Données techniques : identifiant de notification push (OneSignal), jeton d'accès au calendrier, journaux d'activité et adresse IP.</li>
</ul>

<h3>3. Finalités et bases légales</h3>
<ul>
<li>Fourniture des services et gestion du compte — exécution du contrat (CGU) ;</li>
<li>Paiements et facturation — exécution du contrat et obligation légale ;</li>
<li>Notifications et communications — consentement et/ou intérêt légitime ;</li>
<li>Sécurité, prévention des abus et journalisation — intérêt légitime ;</li>
<li>Respect des obligations légales.</li>
</ul>

<h3>4. Destinataires et sous-traitants</h3>
<p>Les données peuvent être communiquées aux prestataires techniques suivants, agissant en qualité de sous-traitants : <strong>OVH</strong> (hébergement), <strong>Stripe</strong> (paiements), <strong>OneSignal</strong> (notifications push), <strong>Google</strong> et <strong>Apple</strong> (authentification), ainsi que le service de messagerie d'OVH (envoi d'e-mails). Aucune donnée n'est vendue à des tiers.</p>

<h3>5. Transferts hors Union européenne</h3>
<p>Certains sous-traitants (notamment Stripe et OneSignal) peuvent traiter des données en dehors de l'Union européenne. Ces transferts sont encadrés par des garanties appropriées, telles que les clauses contractuelles types de la Commission européenne.</p>

<h3>6. Durées de conservation</h3>
<p>Les données sont conservées pendant la durée d'utilisation du compte, puis supprimées ou anonymisées dans un délai raisonnable après sa clôture. Les données de facturation et les journaux d'activité sont conservés conformément aux durées légales applicables.</p>

<h3>7. Sécurité</h3>
<p>UpcycleConnect met en œuvre des mesures techniques et organisationnelles raisonnables pour protéger les données : chiffrement des mots de passe, accès restreint aux données et communications sécurisées.</p>

<h3>8. Cookies et stockage local</h3>
<p>La Plateforme utilise le stockage local du navigateur (localStorage) pour conserver votre jeton de session et vos préférences. Ces éléments sont strictement nécessaires au fonctionnement du service. Le service OneSignal peut déposer des identifiants techniques destinés à l'envoi de notifications, soumis à votre consentement. Aucun cookie publicitaire n'est utilisé.</p>

<h3>9. Vos droits</h3>
<p>Conformément au RGPD, vous disposez des droits d'accès, de rectification, d'effacement, de limitation, d'opposition et de portabilité de vos données, ainsi que du droit de retirer votre consentement à tout moment. Vous pouvez exercer ces droits en écrivant à dpo@upcycleconnect.xyz. Vous disposez également du droit d'introduire une réclamation auprès de la CNIL (https://www.cnil.fr).</p>

<h3>10. Mineurs</h3>
<p>La Plateforme n'est pas destinée aux enfants de moins de 15 ans. Aucune collecte de données les concernant n'est effectuée sciemment.</p>

<h3>11. Modifications</h3>
<p>La présente politique peut être mise à jour. La date de dernière mise à jour est indiquée en haut de cette page.</p>`

const legalAccessibilityContent = `<h2>Déclaration d'accessibilité</h2>
<p>UpcycleConnect s'engage à rendre son site accessible au plus grand nombre, dans l'esprit du Référentiel Général d'Amélioration de l'Accessibilité (RGAA), dans sa version 4.1. Cette déclaration s'applique au site <strong>https://upcycleconnect.xyz</strong>.</p>

<h3>1. Contexte</h3>
<p>UpcycleConnect est un projet étudiant à but pédagogique (voir <a href="/mentions-legales">Mentions légales</a>). À ce titre, il n'entre pas dans le champ des organismes légalement soumis à l'obligation d'accessibilité prévue par l'article 47 de la loi n° 2005-102 du 11 février 2005. Nous avons néanmoins souhaité prendre en compte les bonnes pratiques d'accessibilité et corriger les défauts identifiables.</p>

<h3>2. État de conformité</h3>
<p>Le site UpcycleConnect est <strong>partiellement conforme</strong> au RGAA version 4.1. Cette appréciation repose sur une auto-évaluation et non sur un audit formel certifié portant sur l'intégralité des 106 critères du référentiel.</p>

<h3>3. Résultats des tests</h3>
<p>Une auto-évaluation a été menée à l'aide de l'outil automatisé axe-core (basé sur les règles WCAG 2.1 niveau AA), sur les principales pages du site public et de l'espace d'administration, en thème clair et sombre. Les anomalies détectées ont été corrigées ; à l'issue de ces corrections, ces contrôles automatisés ne relèvent plus d'erreur sur les pages testées. Les outils automatisés ne couvrant qu'une partie des critères, une conformité totale ne peut être garantie sans audit manuel complet.</p>

<h3>4. Mesures d'accessibilité mises en œuvre</h3>
<ul>
<li>contrastes de couleurs renforcés pour le texte (objectif 4.5:1), en thèmes clair et sombre ;</li>
<li>noms accessibles ajoutés aux boutons à icône, aux menus déroulants et aux champs de formulaire (libellés associés) ;</li>
<li>textes alternatifs sur les images porteuses d'information ;</li>
<li>structure des titres et repères de page (landmarks) corrigés ;</li>
<li>navigation au clavier améliorée (zones défilables focusables) ;</li>
<li>balisage HTML sémantique et indication de la langue des pages ;</li>
<li>conception responsive adaptée au zoom et aux différentes tailles d'écran.</li>
</ul>

<h3>5. Contenus non accessibles</h3>
<p>Malgré ces améliorations, certaines limitations peuvent subsister :</p>
<ul>
<li>l'audit manuel de l'ensemble des 106 critères du RGAA (tests avec lecteurs d'écran, parcours complets) n'a pas été réalisé ;</li>
<li>des composants riches (éditeur de texte, sélecteurs de date, cartes) peuvent ne pas être pleinement utilisables avec toutes les technologies d'assistance ;</li>
<li>certains contenus ajoutés par les utilisateurs peuvent être dépourvus d'alternative textuelle.</li>
</ul>

<h3>6. Établissement de cette déclaration</h3>
<p>Cette déclaration repose sur une auto-évaluation. Technologies utilisées pour la réalisation du site : HTML, CSS et JavaScript (Vue.js).</p>

<h3>7. Retour d'information et contact</h3>
<p>Si vous rencontrez un défaut d'accessibilité vous empêchant d'accéder à un contenu ou à une fonctionnalité, vous pouvez nous le signaler :</p>
<ul>
<li>par e-mail : contact@upcycleconnect.xyz ;</li>
<li>via le formulaire de <a href="/contact">contact</a> du site.</li>
</ul>

<h3>8. Voies de recours</h3>
<p>Cette procédure s'applique aux organismes soumis à l'obligation légale. Bien qu'UpcycleConnect n'y soit pas soumis, nous la mentionnons à titre informatif. Si, après nous avoir signalé un défaut, vous n'obtenez pas de réponse satisfaisante, vous pouvez saisir le Défenseur des droits (https://www.defenseurdesdroits.fr).</p>

<h3>9. En savoir plus</h3>
<p>Pour en savoir plus sur le RGAA : <a href="https://accessibilite.numerique.gouv.fr">https://accessibilite.numerique.gouv.fr</a>.</p>`

var legalDefaults = []models.LegalDocument{
	{Slug: "mentions-legales", Title: "Mentions légales", Content: legalMentionsContent},
	{Slug: "cgu-cgv", Title: "Conditions générales d'utilisation et de vente", Content: legalCGUContent},
	{Slug: "confidentialite", Title: "Politique de confidentialité", Content: legalPrivacyContent},
	{Slug: "accessibilite", Title: "Déclaration d'accessibilité", Content: legalAccessibilityContent},
}

var legalReplaceable = map[string]bool{
	"":  true,
	`<h2>Mentions légales</h2><p>Éditeur du site : UpcycleConnect.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>`:                                              true,
	`<h2>Conditions générales</h2><p>Les présentes conditions régissent l'utilisation de la plateforme UpcycleConnect.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>`: true,
	`<h2>Politique de confidentialité</h2><p>UpcycleConnect s'engage à protéger les données personnelles de ses utilisateurs.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>`: true,
}

func SeedLegalDocuments(db *gorm.DB) {
	for _, d := range legalDefaults {
		var existing models.LegalDocument
		err := db.Where("slug = ?", d.Slug).First(&existing).Error
		if err != nil {
			doc := d
			db.Create(&doc)
			continue
		}
		if legalReplaceable[strings.TrimSpace(existing.Content)] {
			existing.Title = d.Title
			existing.Content = d.Content
			db.Save(&existing)
		}
	}
}
