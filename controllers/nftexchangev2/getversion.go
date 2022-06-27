package nftexchangev2

func (nft *NftExchangeControllerV2) GetVersion() {
	var version string
	version = `
			<h1>0.5.18</h1>
			<p>Increase the switch of whether nft upload needs to be reviewed		[Y]</p>
			<h1>0.5.17</h1>
			<p>login Add the approve_addr field to increase the answering mechanism, 
					see the description for details.								[Y]</p>
			<h1>0.5.16.1</h1>
			<p>fix bugs for querynftlist()											[Y]</p>
			<p>remove validation for search()										[Y]</p>
			<h1>0.5.16</h1>
			<p>[new] buyResultInterface												[Y]</p>
			<h1>0.5.15</h1>
			<p>[update] queryHomePage Add field total								[Y]</p>
			<p>[update] modifyCollectionsImage Input interface update				[Y]</p>
			<h1>0.5.14</h1>
			<p>[update] queryCollectionInfo Add 3 transaction fields
			<p>- trade_amount - trade_avg_price - trade_floor_price					[Y]</p>
			<h1>0.5.13</h1>
			<p>[update] queryUserBidlist Add name field NFT name					[Y]</p>
			<p>[update] queryUserOfferlist Add name field NFT name					[Y]</p>
			<h1>0.5.12</h1>
			<p>[update] All addresses are converted to lowercase (user address, contract address)	[Y]</p>
			<h1>0.5.11</h1>
			<p>[update] queryNFT add field collection_desc							[Y]</p>
			<h1>0.5.10</h1>
			<p>[new] version Add version interface									[Y]</p>
			<p>[new] set_sys_para Back-end interface, front-end ignore				[Y]</p>
			<p>[new] get_sys_para Back-end interface, front-end ignore				[Y]</p>
			<h1>0.5.9.1</h1>
			<p>[update] queryUserNFTList add field collection_creator_addr			[Y]</p>
			<p>[update] queryUserCollectionList add field collection_creator_addr	[Y]</p>
			<p>[update] queryUserFavoriteList add field collection_creator_addr		[Y]</p>
			<p>[update] queryNFTCollectionList add field collection_creator_addr	[Y]</p>
			<p>[update] queryCollectionInfo add field collection_creator_addr		[Y]</p>
			<p>[update] queryHomePage add field collection_creator_addr				[Y]</p>
			<p>[update] queryNFT add field trade_hash								[Y]</p>
			<p>[update] queryUserTradingHistroy add field trade_hash				[Y]</p>
			<p>[update] queryMarketTradingHistroy add field trade_hash				[Y]</p>
			`

	nft.Ctx.ResponseWriter.Write([]byte(version))
}
