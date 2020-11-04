$(function(){
    $('#addUsedItemDialog input:not([type="radio"]):not([type="submit"])').css({
        'width': '100%',
        'clear': 'both',
        'margin': '2px 0px'
    });
    $('#addUsedItemDialog .timing').css({
        'width': '100%',
        'clear': 'both',
        'margin': '5px 0px'
    });

    $('#addUsedItemDialog [type="submit"]').css({
        'padding': '10px',
        'border-radius' : '5px',
        'border' : 'none',
        'background-color': '#4472C4',
        'color': '#ffffff',
        'font-size': '12px',
        'line-height': '14px',
        'width' : '20%',
        'margin': '10px 0 0 auto'
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
  })