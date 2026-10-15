#include <stdio.h>
#include <stdlib.h>
typedef struct Node
{
    int data;
    struct Node *left;
    struct Node *right;

} Node;

Node *createNode(int value)
{
    Node *newNode = malloc(sizeof(Node));
    if (newNode == NULL)
    {
        fprintf(stderr, "memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    newNode->data = value;
    newNode->left = NULL;
    newNode->right = NULL;
    return newNode;
}
void freeTree(Node *root)
{
    if (root == NULL)
        return;

    freeTree(root->left);
    freeTree(root->right);
    free(root);
}
void preorder(Node *root)
{
    if (root == NULL)
        return;
    printf("%d\n", root->data);
    preorder(root->left);
    preorder(root->right);
}
void inorder(Node *root)
{
    if (root == NULL)
        return;
    inorder(root->left);
    printf("%d\n", root->data);
    inorder(root->right);
}
void postorder(Node *root)
{
    if (root == NULL)
        return;
    postorder(root->left);
    postorder(root->right);
    printf("%d\n", root->data);
}

int treeHeight(Node *root)
{
    if (root == NULL)
        return 0;
    int leftHeight = 0;
    int rightHeight = 0;
    leftHeight = treeHeight(root->left);
    rightHeight = treeHeight(root->right);
    if (rightHeight > leftHeight)
    {
        return rightHeight + 1;
    }
    return leftHeight + 1;
}
void levelOrder(Node *root)
{
    int capacity = 8;
    if (root == NULL)
        return;

    Node **queue = malloc(capacity * sizeof(Node *));
    if (queue == NULL)
    {
        fprintf(stderr, "queue allocation failed\n");
        exit(EXIT_FAILURE);
    }
    int front = 0;
    int rear = 0;

    queue[rear] = root;
    rear = (rear + 1) % capacity;

    while (front != rear)
    {
        Node *current = queue[front];
        front = (front + 1) % capacity;
        printf("%d\n", current->data);
        if (current->left)
        {
            if ((rear + 1) % capacity == front)
            {
                int newcapacity = capacity * 2;
                int oldCount = capacity - 1;
                Node **newQueue = malloc(newcapacity * sizeof(Node *));
                if (newQueue == NULL)
                {
                    fprintf(stderr, "allocation failed\n");
                    exit(EXIT_FAILURE);
                }
                for(int i=0; i<oldCount;i++){
                    newQueue[i]=queue[(front+i)%capacity];
                }
                free(queue);
                queue=newQueue;
                front=0;
                rear=oldCount;
                capacity=newcapacity;
            }
            queue[rear] = current->left;
            rear = (rear + 1) % capacity;
        }
        if (current->right)
        {
            if ((rear + 1) % capacity == front)
            {
                int newcapacity = capacity * 2;
                int oldCount = capacity - 1;
                Node **newQueue = malloc(newcapacity * sizeof(Node *));
                if (newQueue == NULL)
                {
                    fprintf(stderr, "allocation failed\n");
                    exit(EXIT_FAILURE);
                }
                for(int i=0; i<oldCount;i++){
                    newQueue[i]=queue[(front+i)%capacity];
                }
                free(queue);
                queue=newQueue;
                front=0;
                rear=oldCount;
                capacity=newcapacity;
            }
            queue[rear] = current->right;
            rear = (rear + 1) % capacity;
        }
    }
    free(queue);
}
int countNodes(Node *root)
{
    if (root == NULL)
        return 0;
    int leftCount = countNodes(root->left);
    int rightCount = countNodes(root->right);
    return leftCount + rightCount + 1;
}
int main(void)
{
    Node *root = createNode(0);
    if (root == NULL)
        return 0;
    Node *nodes[255];
    nodes[0] = root;

    for (int i = 1; i < 255; i++)
    {
        nodes[i] = createNode(i);
    }
    for (int i = 0; i < 127; i++)
    {
        nodes[i]->left = nodes[2 * i + 1];
        nodes[i]->right = nodes[2 * i + 2];
    }

    levelOrder(root);
    printf("Tree height: %d\n", treeHeight(root));
    printf("Node count: %d\n", countNodes(root));
    freeTree(root);

    return 0;
}