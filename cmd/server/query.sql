SELECT Books.Authors, Books.Title, JSON_EXTRACT(Tags.Val, '$.text')
FROM Tags
         JOIN Items on Tags.ItemID = Items.OID
         JOIN Books on Items.ParentID = Books.OID
WHERE Tags.TagID = 104
  AND JSON_EXTRACT(Tags.Val, '$.text') <> 'Bookmark'
ORDER BY Items.TimeAlt;