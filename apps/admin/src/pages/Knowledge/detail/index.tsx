import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Button, Space, Spin, Tag } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { getDoc } from '@/services/knowledge';

function normalizeTags(tags?: string | string[]) {
  if (Array.isArray(tags)) {
    return tags.map((tag) => tag.trim()).filter(Boolean);
  }
  if (typeof tags === 'string') {
    return tags
      .split(',')
      .map((tag) => tag.trim())
      .filter(Boolean);
  }
  return [];
}

const KnowledgeDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [doc, setDoc] = useState<API.KnowledgeDoc | null>(null);

  useEffect(() => {
    const fetchDoc = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const result = await getDoc(Number(id));
        if (result) {
          setDoc(result);
        }
      } catch (error) {
        console.error('获取文档详情失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchDoc();
  }, [id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (!doc) {
    return (
      <div style={{ textAlign: 'center', padding: 80, color: '#999' }}>
        文档不存在或加载失败
        <br />
        <Button style={{ marginTop: 16 }} onClick={goBack}>
          返回
        </Button>
      </div>
    );
  }

  return (
    <div>
      <ProCard
        title="文档详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
            <Button>编辑</Button>
            <Button type="primary">发布</Button>
          </Space>
        }
      >
        <ProDescriptions
          column={2}
          dataSource={doc}
          columns={[
            { title: '文档ID', dataIndex: 'id' },
            { title: '标题', dataIndex: 'title' },
            { title: '分类', dataIndex: 'category' },
            {
              title: '公开',
              dataIndex: 'is_public',
              render: (_, record) => (
                <Tag color={record.is_public ? 'green' : 'default'}>
                  {record.is_public ? '公开' : '内部'}
                </Tag>
              ),
            },
            {
              title: '标签',
              dataIndex: 'tags',
              render: (_, record) => {
                const tags = normalizeTags(record.tags);
                return tags.length > 0 ? tags.map((tag) => <Tag key={tag}>{tag}</Tag>) : '-';
              },
            },
            { title: '创建时间', dataIndex: 'created_at' },
            { title: '更新时间', dataIndex: 'updated_at' },
          ]}
        />
      </ProCard>

      <ProCard title="文档内容" style={{ marginTop: 16 }}>
        {doc.content ? (
          <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>
            {doc.content}
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
            暂无文档内容
          </div>
        )}
      </ProCard>
    </div>
  );
};

export default KnowledgeDetailPage;
