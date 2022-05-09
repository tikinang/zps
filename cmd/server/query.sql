SELECT Books.Authors, Books.Title, JSON_EXTRACT(Tags.Val, '$.text'), JSON_EXTRACT(Tags.Val, '$.begin'), Tags.TimeEdt
FROM Tags
         JOIN Items on Tags.ItemID = Items.OID
         JOIN Books on Items.ParentID = Books.OID
WHERE Tags.TagID = 104
  AND JSON_EXTRACT(Tags.Val, '$.text') <> 'Bookmark'
ORDER BY Tags.TimeEdt;