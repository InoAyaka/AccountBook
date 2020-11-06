$(function(){
    // --------- modal addUsedItemDialog ---------
    $('#addUsedItemDialog input:not([type="radio"]):not([type="submit"])').css({
        'width': '100%',
        'clear': 'both',
        'margin': '0px 0px 5px 0px'
    });
    $('#addUsedItemDialog .timing').css({
        'width': '100%',
        'clear': 'both',
        'margin': '0px 0px 5px 0px'
    });

    $('#addUsedItemDialog [type="submit"]').css({
        'padding': '10px',
        'border-radius' : '5px',
        'border' : 'none',
        'background-color': '#4472C4',
        'color': '#ffffff',
        'font-size': '12px',
        'line-height': '14px',
        'width' : '80px',
        'margin': '10px 0 5px auto'
    });

    // 画面表示時にダイアログが表示されないよう設定
    $("#addUsedItemDialog").dialog({
        autoOpen: false,
        modal:true,
        width:400
    });
   
    // ボタンのクリックイベント
    $(".addUsedItemBtn").click(function(){
        //購入品IDを補足
        PurchaseItemID = $(this).data('id');
        $('#addUsedItemDialog input[name="purchaseItemID"]').val(PurchaseItemID);
        //残個数を個数の最大値に設定
        RemainingQuantity = $(this).data('requantity');
        $('#addUsedItemDialog input[name="quantity"]').attr("max",RemainingQuantity);
        
        // ダイアログを表示する
        $("#addUsedItemDialog").dialog("open");
    });

    // --------- modal updateUsedItemDialog ---------
    $('#updateUsedItemDialog input:not([type="checkbox"]):not([type="submit"])').css({
        'width': '100%',
        'clear': 'both',
        'margin': '0px 0px 5px 0px'
    });

    $('#updateUsedItemDialog select').css({
        'margin': '0px 0px 5px 0px'
    })

    $('#updateUsedItemDialog [type="submit"]').css({
        'padding': '10px',
        'border-radius' : '5px',
        'border' : 'none',
        'background-color': '#4472C4',
        'color': '#ffffff',
        'font-size': '12px',
        'line-height': '14px',
        'width' : '80px',
        'margin': '10px 0 5px auto'
    });

    // 画面表示時にダイアログが表示されないよう設定
    $("#updateUsedItemDialog").dialog({
        autoOpen: false,
        modal:true,
        width:650
    });

    $(".updatePurchasedItemBtn").click(function(){
        //変更前の情報取得
        PurchaseItemID = $(this).data('id');
        Month = $(this).data('month');
        PurchasedOn = $(this).data('purchased-on').slice(0,10);
        Shop = $(this).data('shop');
        
        Category = $(this).data('category');
        ItemName = $(this).data('item-name');
        Quantity = $(this).data('quantity');
        PriceIncludingTax = $(this).data('price-including-tax')

        $('#updateUsedItemDialog input[name="purchaseItemID"]').val(PurchaseItemID);
        $('#updateUsedItemDialog #exMonth').text(Month);
        $('#updateUsedItemDialog #exPurchasedOn').text(PurchasedOn);
        $('#updateUsedItemDialog #exShop').text(Shop);
        $('#updateUsedItemDialog #exCategory').text(Category);
        $('#updateUsedItemDialog #exItemName').text(ItemName);
        $('#updateUsedItemDialog #exQuantity').text(Quantity);
        $('#updateUsedItemDialog #exIncTax').text("¥ " + PriceIncludingTax);
        
        // ダイアログを表示する
        $("#updateUsedItemDialog").dialog("open");
    });
  })