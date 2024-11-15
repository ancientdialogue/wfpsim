package shortcut

var MonsterNameToID = map[string]int{
	"hilichurl":                            21010101,
	"woodenshieldwallmitachurl":            21020101,
	"blazingaxemitachurl":                  21020201,
	"rockshieldwallmitachurl":              21020301,
	"frostarmlawachurl":                    21020401,
	"stonehidelawachurl":                   21020501,
	"iceshieldwallmitachurl":               21020601,
	"cracklingaxemitachurl":                21020701,
	"thunderhelmlawachurl":                 21020801,
	"dendroslime":                          20010101,
	"largedendroslime":                     20010201,
	"anemoslime":                           20010301,
	"largeanemoslime":                      20010401,
	"stormterror":                          29010102,
	"lupusboreasdominatorofwolves":         29020101,
	"childe":                               29030103,
	"azhdaha":                              29040101,
	"lasignora":                            29050102,
	"magatsumitakenarukaminomikoto":        29060202,
	"shoukinokamitheprodigal":              29070104,
	"guardianofapepsoasis":                 29080102,
	"alldevouringnarwhal":                  29090101,
	"theknave":                             29100103,
	"eyeofthestorm":                        20020101,
	"electrohypostasis":                    20040101,
	"anemohypostasis":                      20040201,
	"geohypostasis":                        20040301,
	"hydrohypostasis":                      20040401,
	"cryohypostasis":                       20040501,
	"pyrohypostasis":                       20040601,
	"dendrohypostasis":                     20040701,
	"electroslime":                         20010501,
	"largeelectroslime":                    20010601,
	"mutantelectroslime":                   20010701,
	"cryoslime":                            20010801,
	"largecryoslime":                       20010901,
	"hydroslime":                           20011001,
	"largehydroslime":                      20011101,
	"hydrohypostasissummon":                20011601,
	"oceanid":                              20050101,
	"hydromimicboar":                       20050201,
	"hydromimiccrane":                      20050301,
	"hydromimiccrab":                       20050401,
	"hydromimicfinch":                      20050501,
	"hydromimicmallard":                    20050601,
	"hydromimicsquirrel":                   20050701,
	"hydromimicfrog":                       20050801,
	"hydromimicraptor":                     20050901,
	"taintedwatersplittingphantasm":        20051001,
	"veterantaintedwatersplittingphantasm": 20051002,
	"taintedwaterspoutingphantasm":         20051101,
	"veterantaintedwaterspoutingphantasm":  20051102,
	"hydrotulpa":                           20051201,
	"halftulpa":                            20051301,
	"thundermanifestation":                 20070101,
	"hilichurlfighter":                     21010201,
	"woodenshieldhilichurlguard":           21010301,
	"hilichurlshooter":                     21010401,
	"pyrohilichurlshooter":                 21010501,
	"hilichurlgrenadier":                   21010601,
	"hilichurlberserker":                   21010701,
	"cryohilichurlshooter":                 21010901,
	"electrohilichurlshooter":              21011001,
	"rockshieldhilichurlguard":             21011201,
	"cryohilichurlgrenadier":               21011301,
	"iceshieldhilichurlguard":              21011401,
	"unusualhilichurl":                     21011501,
	"electrohilichurlgrenadier":            21011601,
	"pyroslime":                            20011201,
	"largepyroslime":                       20011301,
	"geoslime":                             20011401,
	"largegeoslime":                        20011501,
	"pyroabyssmage":                        22010101,
	"cryoabyssmage":                        22010201,
	"hydroabyssmage":                       22010301,
	"electroabyssmage":                     22010401,
	"abyssheraldwickedtorrents":            22020101,
	"abyssheraldfrostfall":                 22020201,
	"abysslectorvioletlightning":           22030101,
	"abysslectorfathomlessflames":          22030201,
	"rockfondrifthoundwhelp":               22040101,
	"thundercravenrifthoundwhelp":          22040201,
	"rockfondrifthound":                    22050101,
	"thundercravenrifthound":               22050201,
	"goldenwolflord":                       22060101,
	"shadowyhuskstandardbearer":            22070101,
	"shadowyhusklinebreaker":               22070201,
	"shadowyhuskdefender":                  22070301,
	"blackserpentknightwindcutter":         22080101,
	"blackserpentknightrockbreakerax":      22080201,
	"iniquitousbaptist":                    22100101,
	"largeshatterstonebreacherprimus":      22110101,
	"shatterstonebreacherprimus":           22110201,
	"largeovergrownbreacherprimus":         22110301,
	"overgrownbreacherprimus":              22110402,
	"tenebrouspapillatypei":                22120101,
	"tenebrousmimiflora":                   22120201,
	"hydrosamachurl":                       21030101,
	"dendrosamachurl":                      21030201,
	"anemosamachurl":                       21030301,
	"geosamachurl":                         21030401,
	"cryosamachurl":                        21030501,
	"electrosamachurl":                     21030601,
	"anemohilichurlrogue":                  21040101,
	"hydrohilichurlrogue":                  21040201,
	"ruinguard":                            24010101,
	"ruinhunter":                           24010201,
	"ruingrader":                           24010301,
	"ruinserpent":                          24010401,
	"ruincruiser":                          24020101,
	"ruindestroyer":                        24020201,
	"ruindefender":                         24020301,
	"ruinscout":                            24020401,
	"perpetualmechanicalarray":             24021101,
	"ruindrakeskywatch":                    24030101,
	"ruindrakeearthguard":                  24030201,
	"aeonblightdrake":                      24030301,
	"primalconstructrepulsor":              24040101,
	"primalconstructreshaper":              24040201,
	"primalconstructprospector":            24040301,
	"algorithmofsemiintransientmatrixofoverseernetwork": 24050101,
	"reconlogmek":                                   24060101,
	"arithmeticenhancermek":                         24060201,
	"underwatersurveymek":                           24060301,
	"geologicalsurveymek":                           24060401,
	"underwaterpatrolmek":                           24060501,
	"deepwaterassaultmek":                           24060601,
	"nimbleharvestermek":                            24060701,
	"assaultspecialistmek":                          24060801,
	"suppressionspecialistmek":                      24060901,
	"annihilationspecialistmek":                     24061001,
	"areaalertmek":                                  24061101,
	"constructionspecialistmek":                     24061201,
	"veteranarithmeticenhancermek":                  24061301,
	"icewindsuite":                                  24070101,
	"experimentalfieldgenerator":                    24080101,
	"secretsourceautomatonhunterseeker":             24090101,
	"secretsourceautomatonconfigurationdevice":      24090201,
	"basicdefensivepyrotower":                       24810101,
	"advanceddefensivepyrotower":                    24810201,
	"basicdefensiveelectrotower":                    24810301,
	"advanceddefensiveelectrotower":                 24810401,
	"basicdefensivehydrotower":                      24810501,
	"advanceddefensivehydrotower":                   24810601,
	"basicdefensivecryotower":                       24810701,
	"advanceddefensivecryotower":                    24810801,
	"basicdefensiveanemotower":                      24810901,
	"advanceddefensiveanemotower":                   24811001,
	"defensivephysicaltower":                        24811101,
	"automatonsmallslimeballoon":                    24920101,
	"automatonlargeslimeballoon":                    24920201,
	"automatonfungusballoon":                        24920301,
	"veteranfatuiskirmishercryogunnerlegionnaire":   23011101,
	"veteranfatuiskirmisherelectrohammervanguard":   23011301,
	"veteranfatuiskirmisherpyroslingerbracer":       23011601,
	"fatuiskirmishercryogunnerlegionnaire":          23010101,
	"fatuiskirmisherhydrogunnerlegionnaire":         23010201,
	"fatuiskirmisherelectrohammervanguard":          23010301,
	"fatuiskirmishergeochanterbracer":               23010401,
	"fatuiskirmisheranemoboxervanguard":             23010501,
	"fatuiskirmisherpyroslingerbracer":              23010601,
	"fatuipyroagent":                                23020101,
	"veteranfatuipyroagent":                         23021101,
	"fatuielectrocicinmage":                         23030101,
	"fatuicryocicinmage":                            23040101,
	"mirrormaiden":                                  23050101,
	"frostoperative":                                23060101,
	"windoperative":                                 23060201,
	"veteranfrostoperative":                         23060301,
	"veteranwindoperative":                          23060401,
	"treasurehoardersliuliu":                        25010101,
	"treasurehoardersscout":                         25010201,
	"treasurehoarderspyropotioneer":                 25010301,
	"treasurehoardershydropotioneer":                25010401,
	"treasurehoarderselectropotioneer":              25010501,
	"treasurehoarderscryopotioneer":                 25010601,
	"treasurehoardershandyman":                      25010701,
	"treasurehoardersraptor":                        25020101,
	"treasurehoardersmarksman":                      25020201,
	"treasurehoarderscarmen":                        25030101,
	"treasurehoardersgravedigger":                   25030201,
	"treasurehoardersseaman":                        25030301,
	"treasurehoardersboss":                          25040101,
	"millelithsoldier":                              25050101,
	"millelithsergeant":                             25050201,
	"shogunateinfantry":                             25050301,
	"shogunateinfantrycaptain":                      25050401,
	"sangonomiyacohort":                             25050501,
	"treasurehoarderspugilist":                      25060101,
	"treasurehoarderscrusher":                       25070101,
	"treasurehoarderswallenstein":                   25070201,
	"nobushijintouban":                              25080101,
	"nobushihitsukeban":                             25080201,
	"nobushikikouban":                               25080301,
	"yorikisamurai":                                 25080401,
	"maguukenki":                                    25090101,
	"maguukenkilonegale":                            25090201,
	"maguukenkigallopingfrost":                      25090301,
	"maguukenkimaskofterror":                        25090401,
	"kairagifierymight":                             25100202,
	"kairagidancingthunder":                         25100101,
	"ochimushaensorcelledthunder":                   25100301,
	"ochimushacankeredflame":                        25100401,
	"spiritofomenstonedrake":                        25110101,
	"spiritofomenwindglede":                         25110201,
	"spiritofomenpyroscorpion":                      25110301,
	"spiritofomendendrospiritserpent":               25110401,
	"eremiteaxevanguard":                            25210101,
	"eremitecrossbow":                               25210201,
	"eremiteravenbeakhalberdier":                    25210301,
	"eremitelinebreaker":                            25210401,
	"eremitesworddancer":                            25210501,
	"eremitedesertclearwater":                       25310101,
	"eremitesunfrost":                               25310201,
	"eremitedaythunder":                             25310301,
	"eremitestoneenchanter":                         25410101,
	"eremitegalehunter":                             25410201,
	"eremitescorchingloremaster":                    25410301,
	"eremitefloralringdancer":                       25410401,
	"flowingbladeskirmisher":                        25500101,
	"wingedwarriorskysoarer":                        25500201,
	"forgedsandinterrogator":                        25501101,
	"swiftstepseeker":                               25501201,
	"wingedwarriorfeathershedder":                   25501301,
	"spiritspeakerwarriordeterminator":              25501401,
	"swiftstepstormscout":                           25502101,
	"flowingbladegreasesaw":                         25502201,
	"spiritspeakerwarriorpsychokinetic":             25502301,
	"swiftsteparmedcourier":                         25503101,
	"spiritspeakerwarriorcontemplator":              25503201,
	"forgedsandshieldbreaker":                       25504101,
	"wingedwarriorskymaster":                        25504201,
	"forgedsandjavelineer":                          25505101,
	"flowingbladeharpoonthrower":                    25505201,
	"tepetlisaurustribechief":                       25550701,
	"tepetlisauruswarriorshardstriker":              25510101,
	"tepetlisauruswarriorrockbreakerblade":          25510201,
	"yumkasauruswarriorwhirlingleaves":              25510301,
	"yumkasauruswarriorflowingskyfire":              25510401,
	"koholasauruswarriorwaveshuttler":               25510501,
	"koholasauruswarriorreefsplitter":               25510601,
	"qucusauruswarriorheartstarhammer":              25510701,
	"qucusauruswarriorblazingsky":                   25510801,
	"iktomisauruswarrioricytidebearer":              25510901,
	"iktomisauruswarriorcryocrystalcannon":          25511001,
	"cryowhopperflower":                             26010101,
	"pyrowhopperflower":                             26010201,
	"electrowhopperflower":                          26010301,
	"cryoregisvine":                                 26020101,
	"pyroregisvine":                                 26020201,
	"electroregisvine":                              26020301,
	"geovishaphatchling":                            26030101,
	"geovishap":                                     26040101,
	"rulerofthechizhangmountainsyiji":               26040201,
	"rulerofthechizhangmountainstianyu":             26040301,
	"primogeovishap":                                26050101,
	"primordialbathysmalvishap":                     26050601,
	"rimebiterbathysmalvishap":                      26050701,
	"bolteaterbathysmalvishap":                      26050801,
	"battlehardenedprimordialbathysmalvishap":       26051201,
	"electrocicin":                                  26060101,
	"hydrocicin":                                    26060201,
	"cryocicin":                                     26060301,
	"beisht":                                        26080101,
	"floatinghydrofungus":                           26090101,
	"floatingdendrofungus":                          26090201,
	"floatinganemofungus":                           26090301,
	"whirlingelectrofungus":                         26090401,
	"whirlingcryofungus":                            26090501,
	"whirlingpyrofungus":                            26090601,
	"stretchypyrofungus":                            26090701,
	"stretchyanemofungus":                           26090801,
	"stretchygeofungus":                             26090901,
	"stretchyelectrofungus":                         26091001,
	"consecratedhornedcrocodile":                    26100101,
	"consecratedfangedbeast":                        26100201,
	"consecratedredvulture":                         26100301,
	"consecratedscorpion":                           26100401,
	"consecratedflyingserpent":                      26100501,
	"jadeplumeterrorshroom":                         26110101,
	"wingeddendroshroom":                            26120101,
	"wingedcryoshroom":                              26120201,
	"groundedhydroshroom":                           26120301,
	"groundedgeoshroom":                             26120401,
	"setekhwenut":                                   26130101,
	"legatusgolem":                                  26140101,
	"praetoriangolem":                               26140201,
	"huntersray":                                    26151001,
	"armoredcrab":                                   26152101,
	"blubberbeast":                                  26153101,
	"hatjellyfish":                                  26154101,
	"bubblyseahorse":                                26155101,
	"balloctopus":                                   26156101,
	"angelicseahare":                                26157101,
	"sternshieldcrab":                               26160101,
	"cherubicseahare":                               26160201,
	"bubblerseahorse":                               26160301,
	"locallegendarmoredcrab":                        26162101,
	"locallegendhuntersray":                         26162201,
	"locallegendblubberbeast":                       26162301,
	"locallegendcherubicseahare":                    26162401,
	"locallegendhatjellyfish":                       26162501,
	"emperoroffireandiron":                          26170101,
	"millennialpearlseahorse":                       26180101,
	"solitarysuanni":                                26190101,
	"xuanwenbeast":                                  26190201,
	"tepetlisaurwhelp":                              26200101,
	"tepetlisaurus":                                 26200201,
	"koholasaurwhelp":                               26210101,
	"koholasaurus":                                  26210201,
	"yumkasaurwhelp":                                26220101,
	"yumkasaurus":                                   26220201,
	"gluttonousyumkasaurmountainking":               26220301,
	"qucusauruschick":                               26230101,
	"qucusaurus":                                    26230201,
	"goldflamequcusaurtyrant":                       26230301,
	"iktomisauruswhelp":                             26240101,
	"iktomisaurus":                                  26240201,
	"battlehardenediktomisaurus":                    26240301,
	"rockcavernouswayobmanifestation":               26260101,
	"battlehardenedrockcavernouswayobmanifestation": 26260401,
	"flowinvertedwayobmanifestation":                26260201,
	"battlehardenedflowinvertedwayobmanifestation":  26260501,
	"foliarswiftwayobmanifestation":                 26260301,
	"battlehardenedfoliarswiftwayobmanifestation":   26260601,
	"battlehardenedburningaflamewayobmanifestation": 26260701,
	"bitingcoldwayobmanifestation":                  26260801,
	"battlehardenedbitingcoldwayobmanifestation":    26261001,
	"burningaflamewayobmanifestation":               26260901,
	"fluidavataroflava":                             26270101,
	"erodingavataroflava":                           26270201,
	"dendrospecter":                                 20060701,
	"hydrospecter":                                  20060101,
	"geospecter":                                    20060201,
	"anemospecter":                                  20060301,
	"cryospecter":                                   20060401,
	"electrospecter":                                20060501,
	"pyrospecter":                                   20060601,
}
